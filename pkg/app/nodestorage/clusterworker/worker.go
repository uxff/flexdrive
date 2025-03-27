package clusterworker

// tobe instead httpworker

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/pingableif"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils"
)

const RegisterTimeoutSec = 12 // 已注册的超时检测
const RegisterIntervalSec = 5 // 作为worker或master注册间隔

const (
	MsgActionFollow      = "cluster.follow"
	MsgActionUpdateNodes = "cluster.updateNodes"
	// MsgActionKickNode    = "cluster.kick"        // @deprecated
	// MsgActionAddNode     = "cluster.add"         // @deprecated
	// MsgActionEraseMaster = "cluster.erasemaster" // @deprecated
)

const (
	ActiveUnset   = 0
	ActiveOnline  = 1
	ActiveOffline = 2 // cannot reach mate
	Deactivated   = 3 // mate has been kicked out of cluster
)

var workerMgrLock = sync.Mutex{}

type Worker struct {
	Id             string // from redis incr? // uniq in cluster, different from other nodes
	ClusterId      string
	ServiceAddr    string // self addr of listen, must be addressable for other nodes
	MasterId       string
	LastRegistered int64 // timestamp
	Active         int   // is active 0=init, not activated; 1=active; 2=offline; 3=deactivated

	// ClusterMembers    map[string]*Worker `json:"-"`
	// ClusterMembersMap表示节点名单，初始化后基本不变，名单的加减将由接口更新名单来完成，而不是运行时随意更改
	ClusterMembersMap *utils.Map[string, *Worker]
	members           string // list of members in string, eg: 127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033
	listVer           string // the time version that when members is set

	clusterAssist *clusterHelper

	quitChan       chan bool
	masterGoneChan chan bool
	//masterChangeChan chan string

	metaData map[string]interface{}

	pingableWorker pingableif.PingableWorker // pointer to GrpcWorker
}

func NewWorker(serviceAddr string, clusterId string) *Worker {
	w := &Worker{}
	w.ServiceAddr = serviceAddr
	w.ClusterId = clusterId

	w.clusterAssist = NewClusterHelper(w.ClusterId)
	// w.ClusterMembers = make(map[string]*Worker, 0)
	w.ClusterMembersMap = &utils.Map[string, *Worker]{}

	w.Id = w.clusterAssist.genMemberHash(w.ServiceAddr)

	w.quitChan = make(chan bool, 0)
	w.masterGoneChan = make(chan bool, 1)
	//w.masterChangeChan = make(chan string, 1) // useful?

	w.metaData = make(map[string]interface{}, 0)

	return w
}

func (w *Worker) SetPingableWorker(pingableWorker pingableif.PingableWorker) error {
	w.pingableWorker = pingableWorker
	return nil
}

func (w *Worker) Start() error {
	log.Debugf("worker %s will start", w.Id)

	go w.KeepRegistered()
	go w.PerformMaster()

	// 等待别的worker注册成功
	time.Sleep(time.Millisecond * 20)

	log.Debugf("waiting mates registered in")
	//log.Debugf("only %d/%d mates has been registered, continuing checking", registeredCount, len(w.ClusterMembers))
	// assure mate is registered
	for {
		// 检查节点是否就位
		registeredCount := 0
		// for mateId := range w.ClusterMembers {
		memberCnt := w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
			if mate.IsActive() {
				registeredCount++
			}
		})

		// 有半数节点就位就可以继续了
		// if registeredCount > len(w.ClusterMembers)/2 {
		if registeredCount*2 >= memberCnt {
			log.Debugf("%d/%d mates has been registered as active, will start electing", registeredCount, memberCnt)
			break
		}

		//log.Debugf("only %d/%d mates has been registered, continuing checking", registeredCount, len(w.ClusterMembers))
		time.Sleep(time.Millisecond * 100)
	}

	w.ElectMaster()

	for {
		// elect
		select {

		// 来自自身监控master
		case <-w.masterGoneChan:
			log.Debugf("will elect when master %s timeout", w.MasterId)

			// 清掉已经注册的master 需要重新注册
			w.MasterId = ""
			// for mateId := range w.ClusterMembers {
			w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
				mate.MasterId = ""
			})

			w.ElectMaster()

			// 来自队友通知要强制选举 需要延迟回复吗？
		//case <-w.masterShift:

		//case newMasterId := <-w.masterChangeChan:
		//	log.Debugf("master changed from:%s to %s", w.MasterId, newMasterId)
		//	w.Follow(newMasterId)

		case _, ok := <-w.quitChan:
			if !ok {
				log.Debugf("quitChan is closed in start while Start()")
			}
			log.Debugf("recv quit signal, than stop Start()")
			return fmt.Errorf("worker(%s) master(%s) will quit", w.Id, w.MasterId)
		}
	}

	return fmt.Errorf("worker(%s) master(%s) will quit", w.Id, w.MasterId)
}

// redis hashkey: /nota/clusterId.clusterSalt = [md5(addr:port/clusterId+salt):{workerInfo}]
func (w *Worker) RegisterToMates() {
	workerMgrLock.Lock()
	defer workerMgrLock.Unlock()

	// 从redis注册id
	w.LastRegistered = time.Now().Unix() // 更新自己的，没用

	wg := &sync.WaitGroup{}

	// for mateId := range w.ClusterMembers {
	w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
		if mateId == w.Id {
			mate.MarkActive() // 自己直接更新为活跃
			return
		}
		wg.Add(1)
		go func(mateId string, mate *Worker) {
			defer wg.Done()
			// ping metaData: {"masterId":"xx","members":"127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033","clusterId":"mycluster","listVer":"xx"}
			res, err := w.pingableWorker.PingTo(mate.ServiceAddr, w.Id, w.buildPingMetaData())
			if err == nil && res != nil {
				// 发出ping成功跟新本地的mate状态；收到ping方更新对方自己的mate的active状态。
				w.RegisterIn(mateId, res.Get("masterId"))
			}
		}(mateId, mate)
	})

	wg.Wait()

	// flush active status, if master is offline, then notice to elect new one
	// for mateId := range w.ClusterMembers {
	w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
		if mate.Id != w.Id && mate.isTimeout() {
			mate.Active = ActiveOffline
		}
		// mate.Active = !mate.isTimeout()
		if mateId == w.MasterId && !mate.IsActive() {
			// master 超时 通知重新选举
			w.MasterId = ""
			w.masterGoneChan <- true
		}
	})

}

func (w *Worker) KeepRegistered() {
	// 保持注册成功
	for {
		select {
		case <-w.quitChan:
			log.Debugf("recv quit signal, than stop KeepRegistered")
			return
		default:
			// register self
			w.RegisterToMates()

			mateDesc := ""
			w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
				mateDesc += fmt.Sprintf("%s->%s:%d; ", mate.Id, mate.MasterId, mate.Active)
			})

			log.Debugf("id:%s has registered to mates(%s), master:%s", w.Id, mateDesc, w.MasterId)

			time.Sleep(time.Second * RegisterIntervalSec)
		}
	}
}

func (w *Worker) Follow(masterId string) error {
	// target := w.ClusterMembers[masterId]
	target, ok := w.ClusterMembersMap.Load(masterId)
	if !ok {
		log.Errorf("master(%s) not exist when to follow, this cannot be happen", masterId)
		return fmt.Errorf("master(%s) not exist when to follow", masterId)
	}

	// check if master follows me
	if target.MasterId != "" && target.MasterId != target.Id {
		// 跟随主人的主人
		//return w.Follow(target.MasterId)
		log.Warnf("master follows others. should follow master(%s)'s master(%s)?", target.Id, target.MasterId)
	}

	w.MasterId = masterId

	log.Errorf("I(%s) followed %s", w.Id, masterId)
	// as same as PerformFollower
	return nil
}

// 选择出不超时的 至少选择出自己
func (w *Worker) VoteAMaster() string {

	allCondidateMateIds := make([]string, 0)
	memberCnt := w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
		if mate.IsActive() {
			allCondidateMateIds = append(allCondidateMateIds, mateId)
		}
	})

	if memberCnt == 0 {
		return w.Id
	}

	if len(allCondidateMateIds) == 0 {
		// must use self
		log.Debugf("vote to self:%s", w.Id)
		return w.Id
	}

	sort.Strings(allCondidateMateIds)

	expectedMasterId := allCondidateMateIds[0]

	log.Debugf("w(%v) voted master:%v in %d mates:%s", w.Id, expectedMasterId, len(allCondidateMateIds), allCondidateMateIds)

	return expectedMasterId
}

// 选举出主节点并跟随。确保调用前，名单中节点都已经是ActiveOnline。
func (w *Worker) ElectMaster() {
	// 抢占式选举 最快选举好的直接广播给别人 让别人无条件服从
	votedMasterId := w.VoteAMaster()

	log.Debugf("%s voted %s as master and will follow", w.Id, votedMasterId)
	w.Follow(votedMasterId)

	//w.VotedMasterId = ""
	w.BroadcastVoted(votedMasterId)

}

func (w *Worker) PerformMaster() {

	log.Debugf("worker %s will perform master, real-time master:%s", w.Id, w.MasterId)

	tick := time.NewTicker(time.Second * RegisterIntervalSec)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if w.MasterId == w.Id {
				// for mateId := range w.ClusterMembers {
				w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
					if mateId != w.Id && mate.MasterId != w.MasterId {
						// if timeout?
						log.Debugf("demand %s to follow me %s", mateId, w.MasterId)
						go w.DemandFollow(mate, w.MasterId)
					}
				})
			}
		case <-w.quitChan:
			log.Debugf("recv quit signal, than stop PerformMaster")
			return
		}
	}
}

// @deprecated replaced by WrapMetaData()
func (w *Worker) ToString() string {
	buf, _ := json.Marshal(w)
	return string(buf)
}

func (w *Worker) DemandFollow(mate *Worker, masterId string) error {

	_, err := w.pingableWorker.MsgTo(mate.ServiceAddr, MsgActionFollow, "", url.Values{"masterId": {masterId}})
	//res := w.MessageTo("follow", mateId, nil) // instead by pingableWorker.MsgTo

	if err != nil {
		log.Debugf("i:%s demand:%s follow:%s error:%v", w.Id, mate.Id, masterId, err)
		return err
	}

	return nil
}

func (w *Worker) Quit() {
	workerMgrLock.Lock()
	defer workerMgrLock.Unlock()

	if w.quitChan != nil {
		close(w.quitChan)
		w.quitChan = nil
		log.Debugf("quit chan is closed")
	} else {
		log.Debugf("quit chan has already closed yet sometimes ago")
	}
}

func (w *Worker) calcTimeout() int64 {
	return time.Now().Unix() - w.LastRegistered
}

func (w *Worker) isTimeout() bool {
	return w.calcTimeout() > RegisterTimeoutSec
}

func (w *Worker) BroadcastVoted(masterId string) {
	// for mateId := range w.ClusterMembers {
	w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
		if mateId == w.Id {
			// 跳过自己
			return
		}
		go w.DemandFollow(mate, masterId)
	})
}

func (w *Worker) BroadcastMembersUpdated() {
	// for mateId := range w.ClusterMembers {
	w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
		if mateId == w.Id {
			// 跳过自己
			return
		}
		go w.pingableWorker.MsgTo(mate.ServiceAddr, MsgActionUpdateNodes, "", w.buildPingMetaData())
	})
}

// Use should be checked strictly. mateServiceAddrs should be a complete list, including self.
// Will not trigger election.
func (w *Worker) UpdateMates(clusterMembers string, listVer string) (memberCnt int) {

	IamIn := false

	mateServiceAddrs := strings.Split(clusterMembers, ",")
	for _, nodeAddr := range mateServiceAddrs {
		nodeId := w.clusterAssist.genMemberHash(nodeAddr)
		mate, exist := w.ClusterMembersMap.Load(nodeId)
		if !exist {
			mate = NewWorker(nodeAddr, w.ClusterId)
			// w.ClusterMembers[mate.Id] = mate
			w.ClusterMembersMap.Store(mate.Id, mate)
			log.Debugf("%s joined my(%s) cluster(%s)", mate.Id, w.Id, w.ClusterId)
		}
		if mate.Id == w.Id {
			IamIn = true
			mate.MarkActive() // at least myself online
		}
	}

	if !IamIn {
		log.Errorf("I(%s,%s) am not in the list(%s). I quit.", w.Id, w.ServiceAddr, clusterMembers)
		w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
			if mateId != w.Id { // not posible
				w.ClusterMembersMap.Delete(mateId)
			}
		})
		return 0
	}

	// check if ClusterMembers does not include node, then delete it
	memberCnt = w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) {
		exist := false
		for _, inputMateAddr := range mateServiceAddrs {
			if mate.ServiceAddr == inputMateAddr {
				exist = true
				break
			}
		}
		if !exist {
			w.ClusterMembersMap.Delete(mateId)
			log.Debugf("%s be kicked out my(%s) cluster(%s)", mate.Id, w.Id, w.ClusterId)
		}
	})

	log.Debugf("mate list updated:%v cnt:%d ver:%s", mateServiceAddrs, memberCnt, listVer)
	w.members = clusterMembers
	w.listVer = listVer // time.Now().Format("20060102T150405")

	// TODO: how about I was not in the list?
	return
}

// Use should be checked strictly
// func (w *Worker) DeleteMate(mateId string) {
// 	w.ClusterMembersMap.Delete(mateId)
// }

func (w *Worker) RegisterIn(mateId string, masterIdOfMate string) {
	// if _, ok := w.ClusterMembers[mateId]; !ok {
	mate, ok := w.ClusterMembersMap.Load(mateId)
	if !ok {
		mates := make([]string, 0)
		w.ClusterMembersMap.RangeAndCount(func(mateId string, mate *Worker) { mates = append(mates, mateId) })
		log.Debugf("when mate %s register in, not exist in my member list:%+v", mateId, mates)
		return
	}

	mate.MarkActive()
	mate.MasterId = masterIdOfMate
}

func (w *Worker) genServeUrl(method string, params url.Values) string {
	u := url.URL{
		Scheme:   "http",
		Host:     w.ServiceAddr,
		Path:     "/" + method,
		RawQuery: params.Encode(),
	}
	return u.String()
}

func (w *Worker) WrapMetaData() string {
	b, _ := json.Marshal(w.metaData)
	return string(b)
}

func (w *Worker) DecodeMetaData(str string) {
	json.Unmarshal([]byte(str), &w.metaData)
}

func (w *Worker) ServePingable() error {

	// addNode is deprecated. use UpdateNodes instead, will fully repleace the existing list, and re-elect new master
	// @IMPORTANT: receiver must compare mateListVer and local listVer, use new one.
	// @param reqParam: {"members":"node1,node2,...","listVer":"20250101.150405"}
	w.pingableWorker.RegisterMsgHandler(MsgActionUpdateNodes, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		members := reqParam.Get("members")
		if members == "" {
			log.Errorf("when handling UpdateNodes, members not give, ignore.")
			return w.buildMsgRes("FAIL", "no members given"), nil
		}

		clusterId := reqParam.Get("clusterId")
		if clusterId != w.ClusterId {
			log.Errorf("when handling UpdateNodes, clusterId(%s) not matched, ignore.", clusterId)
			return w.buildMsgRes("FAIL", "clusterId("+clusterId+") not matched"), nil
		}

		// compare local version, if remote version is newer, then use remote version and update local.
		mateListVer := reqParam.Get("listVer")
		if mateListVer != "" && mateListVer > w.listVer {

			memberCnt := w.UpdateMates(members, mateListVer)
			if memberCnt == 0 {
				log.Errorf("I(%s) join members failed. memmbers:%s clusterId:%s", w.Id, members, w.ClusterId)
				return w.buildMsgRes("FAIL", "join failed. no members. quit."), nil
			}

			w.MarkMateActive(fromId, ActiveOnline) // at least mark fromId active
			// trigger election
			go w.ElectMaster()
		}

		return w.buildMsgRes("OK", ""), nil
	})

	// will follow param masterId. dont return error
	// @param string masterId
	w.pingableWorker.RegisterMsgHandler(MsgActionFollow, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		masterId := reqParam.Get("masterId")
		if w.MasterId == masterId {
			log.Debugf("i(%s) have already follow %s while recv demand follow", w.Id, masterId)
			return w.buildMsgRes("OK", "already"), nil
		}
		// masterWorker, masterExist := w.ClusterMembers[masterId]
		masterWorker, masterExist := w.ClusterMembersMap.Load(masterId)
		if !masterExist {
			log.Debugf("will follow but masterId:" + masterId + " not exist")
			return w.buildMsgRes("FAIL", "masterId not exist"), nil
		}

		// if I am a master, only follow fromId if compared as I'm larger
		if w.Id == w.MasterId {
			if w.Id < masterId {
				log.Warnf("I(%s) do not fucking follow you mate(%s)", w.Id, masterId)
				return w.buildMsgRes("FAIL", "reject to follow because I'm the prime master"), nil
			}
			log.Warnf("I(%s) am a master, but I surrender to follow %s", w.Id, masterId)
		}

		// masterPingRes := w.PingNode(masterId)
		masterPingRes, err := w.pingableWorker.PingTo(masterWorker.ServiceAddr, w.Id, w.buildPingMetaData())
		if err != nil || masterPingRes == nil {
			log.Debugf("will follow(%s) but ping error:%v", err)
			return w.buildMsgRes("FAIL", "ping master fail"), nil
		}

		// masterId = masterPingRes.Get("masterId") // follow master's master
		w.Follow(masterId)

		go w.RegisterToMates()

		return w.buildMsgRes("OK", ""), nil
	})

	// RegisterPong will be triggered if others call PingTO(toMe)
	// ping metaData: {"masterId":"xx","members":"127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033","clusterId":"mycluster","listVer":"xx"}
	w.pingableWorker.RegisterPong(func(fromId, toId, metaData string) (url.Values, error) {
		// w.MarkActive(fromId, ActiveOnline) // at least mark fromId active
		// log.Debugf("receive ping from:%s meta:%s", fromId, metaData)
		mate, ok := w.ClusterMembersMap.Load(fromId)
		if ok {
			mate.MarkActive()

			// parse mate.MasterId
			reqVal, err := url.ParseQuery(metaData)
			if err != nil {
				log.Errorf("parse pong.metaData failed, from:%s meta:%s", fromId, metaData)
			}

			mateClusterId := reqVal.Get("clusterId")
			if mateClusterId != w.ClusterId {
				log.Warnf("receive ping from %s, what hell your clusterId:%s diff from my:%s", fromId, mateClusterId, w.ClusterId)
			}

			mateMasterId := reqVal.Get("masterId")

			mate.MasterId = mateMasterId

			if mateMasterId != w.MasterId {
				log.Warnf("receive ping from %s master %s diff from my master:%s", fromId, mateMasterId, w.MasterId)
				// TODO: how to do?
			}

			mateListVer := reqVal.Get("listVer")
			if mateListVer != w.listVer {
				log.Warnf("receive ping from %s, what hell your listVer:%s diff from my:%s", fromId, mateListVer, w.listVer)
			}

		}
		return nil, nil
	})

	return w.pingableWorker.Serve(w.ServiceAddr)
}

func (w *Worker) GetPingableWorker() pingableif.PingableWorker {
	return w.pingableWorker
}

func (w *Worker) IsActive() bool {
	return w.Active == ActiveOnline
}

func (w *Worker) MarkMateActive(mateId string, mark int) {
	mate, ok := w.ClusterMembersMap.Load(mateId)
	if ok {
		mate.Active = mark
		if mark == ActiveOnline {
			mate.LastRegistered = time.Now().Unix()
		}
	}
}

func (w *Worker) MarkActive() {
	w.Active = ActiveOnline
	w.LastRegistered = time.Now().Unix()
}

func (w *Worker) MarkOffline() {
	w.Active = ActiveOffline
}

func (w *Worker) MarkDeactive() {
	w.Active = Deactivated
}

func (w *Worker) buildPingMetaData() url.Values {
	return url.Values{
		"masterId":  []string{w.MasterId},
		"members":   []string{w.members},
		"clusterId": []string{w.ClusterId},
		"listVer":   []string{w.listVer},
	}
}

func (w *Worker) buildMsgRes(code, msg string) url.Values {
	return url.Values{
		"code":      []string{code},
		"msg":       []string{msg},
		"masterId":  []string{w.MasterId},
		"members":   []string{w.members},
		"clusterId": []string{w.ClusterId},
		"listVer":   []string{w.listVer},
	}
}

func (w *Worker) GenNewListVer() string {
	return time.Now().Format("20060102T150405")
}

package clusterworker

// tobe instead httpworker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/pingableif"
)

const RegisterTimeoutSec = 2  // 已注册的超时检测
const RegisterIntervalSec = 1 // 作为worker或master注册间隔

const (
	MsgActionFollow      = "cluster.follow"
	MsgActionKickNode    = "cluster.kick"
	MsgActionAddNode     = "cluster.add"
	MsgActionEraseMaster = "cluster.erasemaster"
)

type Worker struct {
	Id             string // from redis incr? // uniq in cluster, different from other nodes
	ClusterId      string
	ServiceAddr    string // self addr of listen, must be addressable for other nodes
	MasterId       string
	LastRegistered int64 // timestamp
	Active         bool  // is active

	ClusterMembers map[string]*Worker `json:"-"`

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
	w.ClusterMembers = make(map[string]*Worker, 0)

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
	log.Printf("worker %s will start", w.Id)

	go w.KeepRegistered()
	go w.PerformMaster()

	// 等待别的worker注册成功
	time.Sleep(time.Millisecond * 20)

	log.Printf("waiting mates registered in")
	//log.Printf("only %d/%d mates has been registered, continuing checking", registeredCount, len(w.ClusterMembers))
	// assure mate is registered
	for {
		// 检查节点是否就位
		registeredCount := 0
		for mateId := range w.ClusterMembers {
			if w.ClusterMembers[mateId].Active {
				registeredCount++
			}
		}

		// 有半数节点就位就可以继续了
		if registeredCount > len(w.ClusterMembers)/2 {
			log.Printf("%d/%d mates has been registered", registeredCount, len(w.ClusterMembers))
			break
		}

		//log.Printf("only %d/%d mates has been registered, continuing checking", registeredCount, len(w.ClusterMembers))
		time.Sleep(time.Millisecond * 100)
	}

	// 抢占式选举 最快选举好的直接广播给别人 让别人无条件服从
	masterId := w.FindFollowedMaster()
	if masterId != "" {
		log.Printf("%s will follow %s from existing cluster", w.Id, masterId)
		w.Follow(masterId)
	} else {
		w.ElectMaster()
	}

	for {
		// elect
		select {

		// 来自自身监控master
		case <-w.masterGoneChan:
			log.Printf("will elect when master %s timeout", w.MasterId)

			// 清掉已经注册的master 需要重新注册
			w.MasterId = ""
			for mateId := range w.ClusterMembers {
				w.ClusterMembers[mateId].MasterId = ""
			}

			w.ElectMaster()

			// 来自队友通知要强制选举 需要延迟回复吗？
		//case <-w.masterShift:

		//case newMasterId := <-w.masterChangeChan:
		//	log.Printf("master changed from:%s to %s", w.MasterId, newMasterId)
		//	w.Follow(newMasterId)

		case _, ok := <-w.quitChan:
			if !ok {
				log.Printf("quitChan is closed in start while Start()")
			}
			log.Printf("recv quit signal, than stop Start()")
			return fmt.Errorf("worker(%s) master(%s) will quit", w.Id, w.MasterId)
		}
	}

	return fmt.Errorf("worker(%s) master(%s) will quit", w.Id, w.MasterId)
}

// redis hashkey: /nota/clusterId.clusterSalt = [md5(addr:port/clusterId+salt):{workerInfo}]
func (w *Worker) RegisterToMates() {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	// 从redis注册id
	w.LastRegistered = time.Now().Unix() // 没用

	wg := &sync.WaitGroup{}

	for mateId := range w.ClusterMembers {
		wg.Add(1)
		go func(mateId string) {
			defer wg.Done()
			//w.PingNode(mateId)
			if w.Id != mateId {
				metaData := url.Values{
					"masterId": w.MasterId,
				}
				res, err := w.pingableWorker.PingTo(w.ClusterMembers[mateId].ServiceAddr, w.Id, metaData)
				if err == nil && res != nil {
					w.RegisterIn(mateId, res.Get("masterId"))
				}
			}
		}(mateId)
	}

	wg.Wait()

	// flush active status
	for mateId := range w.ClusterMembers {
		w.ClusterMembers[mateId].Active = !w.ClusterMembers[mateId].isTimeout()
		if mateId == w.MasterId && !w.ClusterMembers[mateId].Active {
			// master 超时 通知重新选举
			w.MasterId = ""
			w.masterGoneChan <- true
		}
	}

}

func (w *Worker) KeepRegistered() {
	// 保持注册成功
	for {
		select {
		case <-w.quitChan:
			log.Printf("recv quit signal, than stop KeepRegistered")
			return
		default:
			// register self
			w.RegisterToMates()

			log.Printf("id:%s has registered to mates, master:%s", w.Id, w.MasterId)

			time.Sleep(time.Second * RegisterIntervalSec)
		}
	}
}

func (w *Worker) Follow(masterId string) {
	target := w.ClusterMembers[masterId]
	if target.MasterId != "" && target.MasterId != target.Id {
		// 跟随主人的主人
		//return w.Follow(target.MasterId)
		log.Printf("will follow master(%s)'s master(%s)?", target.Id, target.MasterId)
	}

	w.MasterId = masterId
	w.RegisterToMates()
	// as same as PerformFollower
}

// 选择出不超时的 至少选择出自己
func (w *Worker) VoteMaster() string {
	if len(w.ClusterMembers) == 0 {
		// must use self
		return w.Id
	}

	allMateIds := make([]string, 0)
	for mateId := range w.ClusterMembers {
		if !w.Active {
			// 超时的节点不能参与投票
			continue
		}
		allMateIds = append(allMateIds, mateId)
	}

	if len(allMateIds) == 0 {
		// must use self
		log.Printf("vote to self:%s", w.Id)
		return w.Id
	}

	sort.Strings(allMateIds)

	expectedMasterId := allMateIds[0]

	log.Printf("w(%v) elected master:%v", w.Id, expectedMasterId)

	return expectedMasterId
}

func (w *Worker) ElectMaster() {
	// 抢占式选举 最快选举好的直接广播给别人 让别人无条件服从
	votedMasterId := w.VoteMaster()
	w.BroadcastVoted(votedMasterId)

	w.Follow(votedMasterId)
	//w.VotedMasterId = ""
}

func (w *Worker) FindFollowedMaster() string {
	masterMap := make(map[string]int, 0)
	for mateId := range w.ClusterMembers {
		if w.ClusterMembers[mateId].Active {
			masterMap[mateId]++
		}
	}

	for masterId, followerNum := range masterMap {
		if followerNum >= len(w.ClusterMembers)/2 {
			return masterId
		}
	}

	return ""
}

//
func (w *Worker) PerformMaster() {

	log.Printf("worker %s will perform master", w.Id)

	tick := time.NewTicker(time.Second * RegisterIntervalSec)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if w.MasterId == w.Id {
				for mateId := range w.ClusterMembers {
					if mateId != w.Id && w.ClusterMembers[mateId].MasterId != w.MasterId {
						// if timeout?
						log.Printf("demand %s to follow me %s", mateId, w.MasterId)
						go w.DemandFollow(mateId, w.MasterId)
					}
				}
			}
		case <-w.quitChan:
			log.Printf("recv quit signal, than stop PerformMaster")
			return
		}
	}
}

// @deprecated replaced by WrapMetaData()
func (w *Worker) ToString() string {
	buf, _ := json.Marshal(w)
	return string(buf)
}

func (w *Worker) DemandFollow(mateId string, masterId string) error {

	_, err := w.pingableWorker.MsgTo(mateId, MsgActionFollow, "", url.Values{"masterId": {masterId}})
	//res := w.MessageTo("follow", mateId, nil) // instead by pingableWorker.MsgTo

	if err != nil {
		log.Printf("i:%s demand:%s follow:%s error:%v", w.Id, mateId, masterId, err)
		return err
	}

	return nil
}

func (w *Worker) Quit() {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	if w.quitChan != nil {
		close(w.quitChan)
		w.quitChan = nil
		log.Printf("quit chan is closed")
	} else {
		log.Printf("quit chan has already closed yet sometimes ago")
	}
}

func (w *Worker) calcTimeout() int64 {
	return time.Now().Unix() - w.LastRegistered
}

func (w *Worker) isTimeout() bool {
	return w.calcTimeout() > RegisterTimeoutSec
}

func (w *Worker) BroadcastVoted(masterId string) {
	for mateId := range w.ClusterMembers {
		if mateId == w.Id {
			// 跳过自己
			continue
		}
		go w.DemandFollow(mateId, masterId)
	}
}

func (w *Worker) AddMates(mateServiceAddrs []string) {
	for _, node := range mateServiceAddrs {
		mate := NewWorker(node, w.ClusterId)
		w.ClusterMembers[mate.Id] = mate
	}
}

func (w *Worker) RegisterIn(mateId string, masterIdOfMate string) {
	if _, ok := w.ClusterMembers[mateId]; !ok {
		log.Printf("when %s register in, not exist, my members:%+v", mateId, w.ClusterMembers)
		return
	}

	w.ClusterMembers[mateId].LastRegistered = time.Now().Unix()
	w.ClusterMembers[mateId].MasterId = masterIdOfMate
	w.ClusterMembers[mateId].Active = true
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

	//w.pingableWorker = NewGrpcWorker(w)
	// todo: Follow,Add,Remove,EraseMaster should use native gRPC functions
	// register is only for outer biz

	// @param string nodes node1,node2
	w.pingableWorker.RegisterMsgHandler(MsgActionAddNode, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		if nodesStr := reqParam.Get("nodes"); nodesStr != "" {
			nodesArr := strings.Split(nodesStr, ",")
			w.AddMates(nodesArr)
		}
		return nil, nil
	})

	// @param string nodeId
	w.pingableWorker.RegisterMsgHandler(MsgActionKickNode, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		delete(w.ClusterMembers, reqParam.Get("nodeId"))
		return nil, nil
	})

	w.pingableWorker.RegisterMsgHandler(MsgActionEraseMaster, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		w.MasterId = ""
		return nil, nil
	})

	// dont return error
	// @param string masterId
	w.pingableWorker.RegisterMsgHandler(MsgActionFollow, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		masterId := reqParam.Get("masterId")
		if w.MasterId == masterId {
			log.Printf("i(%s) have already follow %s while recv demand follow", w.Id, masterId)
			return nil, nil
		}
		masterWorker, masterExist := w.ClusterMembers[masterId]
		if !masterExist {
			log.Printf("will follow but masterId:" + masterId + " not exist")
			return nil, nil
		}

		// masterPingRes := w.PingNode(masterId)
		metaData := url.Values{
			"masterId": w.MasterId,
		}
		masterPingRes, err := w.pingableWorker.PingTo(masterWorker.ServiceAddr, w.Id, metaData))
		if err != nil || masterPingRes == nil {
			log.Printf("will follow(%s) but ping error:%v", err)
			return nil, err
		}
		// if masterPingRes.Code != 0 {
		// 	log.Printf("will follow(%s) but ping error:" + masterPingRes.Get("masterId"))
		// 	return nil, nil
		// }

		masterId = masterPingRes.Get("masterId") // follow master's master
		w.Follow(masterId)
		return nil, nil
	})

	return w.pingableWorker.Serve(w.ServiceAddr)
}

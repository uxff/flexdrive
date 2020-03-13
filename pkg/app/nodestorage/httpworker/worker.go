package httpworker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	//"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/uxff/flexdrive/pkg/log"
)

const RegisterTimeoutSec = 200  // 已注册的超时检测
const RegisterIntervalSec = 100 // 作为worker或master注册间隔

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
	router *gin.Engine
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

	w.router = gin.Default()

	return w
}

func (w *Worker) Start() error {
	log.Debugf("worker %s will start", w.Id)

	go w.KeepRegistered()
	go w.PerformMaster()

	// 等待别的worker注册成功
	time.Sleep(time.Millisecond * 20)

	//log.Debugf("waiting mates registered in")
	log.Debugf("waiting mates registered in")
	//log.Debugf("only %d/%d mates has been registered, continuing checking", registeredCount, len(w.ClusterMembers))
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
			log.Debugf("%d/%d mates has been registered", registeredCount, len(w.ClusterMembers))
			break
		}

		//log.Debugf("only %d/%d mates has been registered, continuing checking", registeredCount, len(w.ClusterMembers))
		time.Sleep(time.Millisecond * 100)
	}

	// 抢占式选举 最快选举好的直接广播给别人 让别人无条件服从
	masterId := w.FindFollowedMaster()
	if masterId != "" {
		log.Debugf("%s will follow %s from existing cluster", w.Id, masterId)
		w.Follow(masterId)
	} else {
		w.ElectMaster()
	}

	for {
		// elect
		select {

		// 来自自身监控master
		case <-w.masterGoneChan:
			log.Debugf("will elect when master %s timeout", w.MasterId)

			// 清掉已经注册的master 需要重新注册
			w.MasterId = ""
			for mateId := range w.ClusterMembers {
				w.ClusterMembers[mateId].MasterId = ""
			}

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
			w.PingNode(mateId)
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
			log.Debugf("recv quit signal, than stop KeepRegistered")
			return
		default:
			// register self
			w.RegisterToMates()

			log.Debugf("id:%s has registered to mates, master:%s", w.Id, w.MasterId)

			time.Sleep(time.Second * RegisterIntervalSec)
		}
	}
}

func (w *Worker) Follow(masterId string) {
	target := w.ClusterMembers[masterId]
	if target.MasterId != "" && target.MasterId != target.Id {
		// 跟随主人的主人
		//return w.Follow(target.MasterId)
		log.Debugf("will follow master(%s)'s master(%s)?", target.Id, target.MasterId)
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
		log.Debugf("vote to self:%s", w.Id)
		return w.Id
	}

	sort.Strings(allMateIds)

	expectedMasterId := allMateIds[0]

	log.Debugf("w(%v) elected master:%v", w.Id, expectedMasterId)

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

// useful?
func (w *Worker) PerformMaster() {

	log.Debugf("worker %s will perform master", w.Id)

	tick := time.NewTicker(time.Second * RegisterIntervalSec)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if w.MasterId == w.Id {
				for mateId := range w.ClusterMembers {
					if mateId != w.Id && w.ClusterMembers[mateId].MasterId != w.MasterId {
						// if timeout?
						log.Debugf("demand %s to follow me %s", mateId, w.MasterId)
						go w.DemandFollow(mateId, w.MasterId)
					}
				}
			}
		case <-w.quitChan:
			log.Debugf("recv quit signal, than stop PerformMaster")
			return
		}
	}
}

// 在命令行主动要求的时候调用
//func (w *Worker) EraseRegisteredMaster() {
//	for _, mate := range w.ClusterMembers {
//		go w.DemandEraseMaster(mate)
//	}
//}

// 通过命令行主动要求删除的时候才调用这里 选举和意外掉线不调用这里
//func (w *Worker) EraseRegisteredWorker(workerId string) {
//	for _, mate := range w.ClusterMembers {
//		go w.DemandRemoveWorker(mate, workerId)
//	}
//}

func (w *Worker) PingNode(workerId string) *PingRes {
	if workerId == w.Id {
		w.RegisterIn(workerId, w.MasterId)
		return &PingRes{Code: 0, WorkerId: w.Id, MasterId: w.MasterId, Members: w.ClusterMembers}
	}

	res := w.MessageTo("ping", workerId, nil)

	if res.Code != 0 {
		log.Debugf("ping failed:%v", res)
		return res
	}

	w.RegisterIn(workerId, res.MasterId)

	// todo 如果收到的mate.MasterId和自己的不一样怎么办？

	return res
}

func (w *Worker) MessageTo(method string, targetId string, val url.Values) *PingRes {
	res := &PingRes{}

	target := w.ClusterMembers[targetId]

	if target == nil {
		res.Msg = fmt.Sprintf("worker(%s) has no target when pingNode(%s)", w.Id, targetId)
		res.Code = 11
		return res
	}

	if val == nil {
		val = make(url.Values)
	}

	val.Set("fromId", w.Id)
	val.Set("masterId", w.MasterId)

	targetUrl := target.genServeUrl(method, val)
	resp, err := http.Get(targetUrl)
	if err != nil {
		res.Msg = "Http Error:" + err.Error()
		res.Code = 13
		return res
	}

	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		res.Msg = "Read Response Error:" + err.Error()
		res.Code = 14
		return res
	}

	err = json.Unmarshal(buf, res)
	if err != nil {
		res.Msg = "Unmarshall Error:" + err.Error()
		res.Code = 15
	}

	return res
}

//
func (w *Worker) ToString() string {
	buf, _ := json.Marshal(w)
	return string(buf)
}

func (w *Worker) DemandFollow(mateId string, masterId string) error {

	res := w.MessageTo("follow", mateId, nil)

	if res.Code != 0 {
		log.Debugf("i:%s demand:%s follow:%s error:%v", w.Id, mateId, masterId, res.Msg)
		return fmt.Errorf(res.Msg)
	}

	return nil
}

//func (w *Worker) DemandEraseMaster(mate *Worker) {
//
//}
//func (w *Worker) DemandRemoveWorker(mate *Worker, workerId string) {
//
//}
//func (w *Worker) DemandCollectVotedMaster(mate *Worker, workerId string) {
//
//}

func (w *Worker) Quit() {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

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
		log.Debugf("when %s register in, not exist, my members:%+v", mateId, w.ClusterMembers)
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

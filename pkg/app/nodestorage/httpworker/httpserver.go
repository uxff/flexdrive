package httpworker

import (
	"encoding/json"

	//"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/log"
)

type PingRes struct {
	Code     int
	Msg      string
	WorkerId string
	MasterId string
	Members  map[string]*Worker

	NodeEnt *dao.Node
}

func (w *Worker) ServePingable() error {

	router := w.router

	router.GET("/", func(c *gin.Context) {
		w.JsonOk(c)
	})

	// 接收对方的ping 表示良好
	router.GET("/ping", func(c *gin.Context) {
		fromId := c.Query("fromId")
		if fromId == "" {
			w.JsonError(c, "fromId must not be empty")
			return
		}

		if _, ok := w.ClusterMembers[fromId]; !ok {
			w.JsonError(c, "fromId:"+fromId+" not exist")
			return
		}

		masterId := c.Query("masterId")

		w.RegisterIn(fromId, masterId)

		w.JsonOk(c)
	})

	// 增加节点 支持批量添加
	// @param nodes=http://127.0.0.1:10010,http://127.0.0.1:10011
	router.GET("/add", func(c *gin.Context) {
		nodesStr := c.Query("nodes")
		if nodesStr == "" {
			w.JsonError(c, "nodes must not be empty")
			return
		}

		nodesArr := strings.Split(nodesStr, ",")
		// todo 通知别人add
		w.AddMates(nodesArr)

		w.JsonOk(c)
	})

	// 删除节点
	router.GET("/remove", func(c *gin.Context) {
		nodeId := c.Query("nodeId")
		if nodeId == "" {
			w.JsonError(c, "nodeId must no be empty")
			return
		}

		if nodeId == w.Id {
			w.Quit()
			w.JsonOk(c)
			return
		}

		delete(w.ClusterMembers, nodeId)
		w.JsonOk(c)
	})

	// 被命令跟随某个master
	router.GET("/follow", func(c *gin.Context) {
		fromId := c.Query("fromId")
		if fromId == "" {
			w.JsonError(c, "fromId must no be empty")
			return
		}

		masterId := c.Query("masterId")
		if masterId == "" {
			w.JsonError(c, "masterId must no be empty")
			return
		}

		if masterId == w.MasterId {
			log.Debugf("i have already follow %s while recv demand follow", masterId)
			w.JsonOk(c)
			return
		}

		if _, ok := w.ClusterMembers[masterId]; !ok {
			w.JsonError(c, "will follow but masterId:"+masterId+" not exist")
			return
		}

		masterPingRes := w.PingNode(masterId)
		if masterPingRes.Code != 0 {
			w.JsonError(c, "will follow(%s) but ping error:"+masterPingRes.Msg)
			return
		}

		masterId = masterPingRes.MasterId

		w.Follow(masterId)
		log.Debugf("%s demand me(%s) follow: %s", fromId, w.Id, masterId)
		w.JsonOk(c)

	})

	// 删除master 重新选举
	router.GET("/erasemaster", func(c *gin.Context) {
		masterId := c.Query("masterId")
		if masterId == "" {
			w.JsonError(c, "node must no be empty")
			return
		}

		w.MasterId = ""
		//w.masterGoneChan <- true
		log.Debugf("erasemaster: %s", masterId)
		w.JsonOk(c)
	})

	// 其他节点向本节点提交其投票
	//router.GET("/collectvotedmaster", func(c *gin.Context) {
	//	fromId := c.Query("fromId")
	//	if fromId == "" {
	//		w.JsonError(c, "fromId must no be empty")
	//		return
	//	}
	//
	//	voteId := c.Query("voteId")
	//	if voteId == "" {
	//		w.JsonError(c, "voteId must no be empty")
	//		return
	//	}
	//
	//	if _, ok := w.ClusterMembers[fromId]; !ok {
	//		w.JsonError(c, "fromId:"+fromId+" not exist")
	//		return
	//	}
	//
	//	w.ClusterMembers[fromId].VotedMasterId = voteId
	//	log.Debugf("collect from %s voted master %s", fromId, voteId)
	//	w.JsonOk(c)
	//})

	return router.Run(w.ServiceAddr)

}

func (w *Worker) AttachPostAction(ruoteUrl string, fn func(c *gin.Context)) {
	w.router.POST(ruoteUrl, fn)
}

func (w *Worker) AttachGetAction(ruoteUrl string, fn func(c *gin.Context)) {
	w.router.GET(ruoteUrl, fn)
}

func (w *Worker) JsonError(c *gin.Context, msg string) {
	c.IndentedJSON(200, PingRes{
		Code:     1,
		Msg:      msg,
		WorkerId: w.Id,
		MasterId: w.MasterId,
		Members:  w.ClusterMembers,
	})
}
func (w *Worker) JsonOk(c *gin.Context) {
	c.IndentedJSON(200, PingRes{
		Code:     0,
		Msg:      "ok",
		WorkerId: w.Id,
		MasterId: w.MasterId,
		Members:  w.ClusterMembers,
	})
}

func newPingRes(buf []byte) *PingRes {
	res := &PingRes{}
	err := json.Unmarshal(buf, res)
	if err != nil {
		res.Msg = "Unmarshall PingRes Error:" + err.Error()
		res.Code = 11
	}
	return res
}

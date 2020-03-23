package grpcworker

import (
	"context"
	"encoding/json"
	"net"

	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/app/nodestorage/grpcworker/pb/pingablepb"
	"github.com/uxff/flexdrive/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PingRes struct {
	Code     int
	Msg      string
	WorkerId string
	MasterId string
	Members  map[string]*Worker
}

const (
	MsgActionFollow      = "cluster.follow"
	MsgActionKickNode    = "cluster.kick"
	MsgActionAddNode     = "cluster.add"
	MsgActionEraseMaster = "cluster.erasemaster"
)

type MsgHandler func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error)

// GrpcWorker implements this interface
type PingableWorker interface {
	//Start() error
	RegisterMsgHandler(action string, handler MsgHandler)
}

type GrpcWorker struct {
	worker *Worker

	msgHandlerMap map[string]MsgHandler
}

func NewGrpcWorker(worker *Worker) *GrpcWorker {
	return &GrpcWorker{
		worker: worker,
	}
}

func (g *GrpcWorker) Ping(ctx context.Context, req *pingablepb.PingRequest) (*pingablepb.PingResponse, error) {

	g.worker.RegisterIn(req.FromId, req.MasterId)
	res := &pingablepb.PingResponse{
		Code: 0,
		Msg:  "ok",
		//Members:  w.Members,
		MetaData: g.worker.WrapMetaData(),
	}
	return res, nil
}

func (g *GrpcWorker) Msg(ctx context.Context, req *pingablepb.MsgRequest) (*pingablepb.MsgResponse, error) {
	res := &pingablepb.MsgResponse{
		MsgId:  req.MsgId,
		Action: req.Action,
	}
	if handler, exist := g.msgHandlerMap[req.Action]; exist {
		reqVal, err := url.ParseQuery(req.Data)
		if err != nil {
			res.Code, res.Msg = 1, err.Error()
			return res, nil
		}
		resVal, err := handler(req.FromId, req.ToId, req.MsgId, reqVal)
		if err != nil {
			res.Code, res.Msg = 1, err.Error()
			return res, nil
		}
		if resVal != nil {
			res.Data = resVal.Encode()
		}
		return res, nil
	}

	log.Debugf("when recv msg, no handler regisgered: %v", req)
	res.Code = 99
	res.Msg = "no handler registered, ignore"

	return res, nil
}

func (g *GrpcWorker) RegisterMsgHandler(action string, handler MsgHandler) {

	if g.msgHandlerMap == nil {
		g.msgHandlerMap = make(map[string]MsgHandler, 0)
	}
	g.msgHandlerMap[action] = handler
}

func (w *Worker) ServePingable() error {

	w.pingableWorker = NewGrpcWorker(w)

	w.pingableWorker.RegisterMsgHandler(MsgActionAddNode, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		if nodesStr := reqParam.Get("nodes"); nodesStr != "" {
			nodesArr := strings.Split(nodesStr, ",")
			w.AddMates(nodesArr)
		}
		return nil, nil
	})

	w.pingableWorker.RegisterMsgHandler(MsgActionKickNode, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		delete(w.ClusterMembers, reqParam.Get("nodeId"))
		return nil, nil
	})

	w.pingableWorker.RegisterMsgHandler(MsgActionEraseMaster, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		w.MasterId = ""
		return nil, nil
	})

	w.pingableWorker.RegisterMsgHandler(MsgActionFollow, func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {
		masterId := reqParam.Get("masterId")
		if w.MasterId == masterId {
			return nil, nil
		}
		if _, ok := w.ClusterMembers[masterId]; !ok {
			return nil, nil
		}
		masterPingRes := w.PingNode(masterId)
		if masterPingRes.Code != 0 {
			// w.jsonError(c, "will follow(%s) but ping error:"+masterPingRes.Msg)
			return nil, nil
		}
		masterId = masterPingRes.MasterId // follow master's master
		w.Follow(masterId)
		return nil, nil
	})

	// 开启RPC服务
	lis, err := net.Listen("tcp", w.ServiceAddr)
	if err != nil {
		log.Errorf("监听端口失败：%v", err)
		return err
	}

	// gRPC通用启动流程
	var opts []grpc.ServerOption
	rpcServer := grpc.NewServer(opts...)

	pingablepb.RegisterPingableInterfaceServer(rpcServer, &GrpcWorker{})
	reflection.Register(rpcServer)

	return rpcServer.Serve(lis)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		w.jsonOk(c)
	})

	// 接收对方的ping 表示良好
	router.GET("/ping", func(c *gin.Context) {
		fromId := c.Query("fromId")
		if fromId == "" {
			w.jsonError(c, "fromId must no be empty")
			return
		}

		if _, ok := w.ClusterMembers[fromId]; !ok {
			w.jsonError(c, "fromId:"+fromId+" not exist")
			return
		}

		masterId := c.Query("masterId")

		w.RegisterIn(fromId, masterId)

		w.jsonOk(c)
	})

	// 增加节点 支持批量添加 使用msg替代
	// @param nodes=http://127.0.0.1:10010,http://127.0.0.1:10011
	router.GET("/add", func(c *gin.Context) {
		nodesStr := c.Query("nodes")
		if nodesStr == "" {
			w.jsonError(c, "nodes must not be empty")
			return
		}

		nodesArr := strings.Split(nodesStr, ",")
		// todo 通知别人add
		w.AddMates(nodesArr)

		w.jsonOk(c)
	})

	// 删除节点 使用msg替代
	router.GET("/remove", func(c *gin.Context) {
		nodeId := c.Query("nodeId")
		if nodeId == "" {
			w.jsonError(c, "nodeId must no be empty")
			return
		}

		if nodeId == w.Id {
			w.Quit()
			w.jsonOk(c)
			return
		}

		delete(w.ClusterMembers, nodeId)
		w.jsonOk(c)
	})

	// 被命令跟随某个master 使用msg替代
	router.GET("/follow", func(c *gin.Context) {
		fromId := c.Query("fromId")
		if fromId == "" {
			w.jsonError(c, "fromId must no be empty")
			return
		}

		masterId := c.Query("masterId")
		if masterId == "" {
			w.jsonError(c, "masterId must no be empty")
			return
		}

		if masterId == w.MasterId {
			log.Errorf("i have already follow %s while recv demand follow", masterId)
			w.jsonOk(c)
			return
		}

		if _, ok := w.ClusterMembers[masterId]; !ok {
			w.jsonError(c, "will follow but masterId:"+masterId+" not exist")
			return
		}

		masterPingRes := w.PingNode(masterId)
		if masterPingRes.Code != 0 {
			w.jsonError(c, "will follow(%s) but ping error:"+masterPingRes.Msg)
			return
		}

		masterId = masterPingRes.MasterId

		w.Follow(masterId)
		log.Debugf("%s demand me(%s) follow: %s", fromId, w.Id, masterId)
		w.jsonOk(c)

	})

	// 删除master 重新选举 是用msg替代
	router.GET("/erasemaster", func(c *gin.Context) {
		masterId := c.Query("masterId")
		if masterId == "" {
			w.jsonError(c, "node must no be empty")
			return
		}

		w.MasterId = ""
		//w.masterGoneChan <- true
		log.Debugf("erasemaster: %s", masterId)
		w.jsonOk(c)
	})

	// 其他节点向本节点提交其投票
	//router.GET("/collectvotedmaster", func(c *gin.Context) {
	//	fromId := c.Query("fromId")
	//	if fromId == "" {
	//		w.jsonError(c, "fromId must no be empty")
	//		return
	//	}
	//
	//	voteId := c.Query("voteId")
	//	if voteId == "" {
	//		w.jsonError(c, "voteId must no be empty")
	//		return
	//	}
	//
	//	if _, ok := w.ClusterMembers[fromId]; !ok {
	//		w.jsonError(c, "fromId:"+fromId+" not exist")
	//		return
	//	}
	//
	//	w.ClusterMembers[fromId].VotedMasterId = voteId
	//	log.Debugf("collect from %s voted master %s", fromId, voteId)
	//	w.jsonOk(c)
	//})

	return router.Run(w.ServiceAddr)

}

func (w *Worker) jsonError(c *gin.Context, msg string) {
	c.IndentedJSON(200, PingRes{
		Code:     1,
		Msg:      msg,
		WorkerId: w.Id,
		MasterId: w.MasterId,
		Members:  w.ClusterMembers,
	})
}
func (w *Worker) jsonOk(c *gin.Context) {
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

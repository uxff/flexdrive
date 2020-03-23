package grpcworker

import (
	"context"
	"errors"
	"net"

	"net/url"
	"strings"

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

// 消息处理句柄 问题：回到弱类型 但是能兼容grpc和http的实现
type MsgHandler func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error)

// 通信组件 每一个既可以收信又可以发信
// GrpcWorker implements this interface
type PingableWorker interface {
	Serve() error
	//Ping()                                                // ping out to other
	//Pong()                                                // like recv ping
	RegisterMsgHandler(action string, handler MsgHandler) // like recv OnMsg
	MsgTo(fromId, toId, action, msgId string, param url.Values) (url.Values, error)
}

// 集成GrpcServer 和 grpc client
type GrpcWorker struct {
	worker *Worker

	rpcServer *grpc.Server

	// need a map of connections
	rpcClientMap map[string]pingablepb.PingableInterfaceClient

	msgHandlerMap map[string]MsgHandler
}

func NewGrpcWorker(worker *Worker) *GrpcWorker {

	return &GrpcWorker{
		worker:       worker,
		rpcClientMap: make(map[string]pingablepb.PingableInterfaceClient),
	}
}

// on ping
func (g *GrpcWorker) Ping(ctx context.Context, req *pingablepb.PingRequest) (*pingablepb.PingResponse, error) {

	g.worker.RegisterIn(req.FromId, req.MasterId)
	res := &pingablepb.PingResponse{
		Code: 0,
		Msg:  "ok",
		//Members:  w.Members,
		MetaData: g.worker.WrapMetaData(), // todo rely on interface
	}
	return res, nil
}

// on msg
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

// msg out
func (g *GrpcWorker) MsgTo(fromId, toId, action, msgId string, param url.Values) (url.Values, error) {
	req := &pingablepb.MsgRequest{
		FromId: fromId,
		ToId:   toId,
		MsgId:  msgId,
		Action: action,
		Data:   param.Encode(),
	}

	ctx := context.Background()
	rpcClient := g.getClient(toId)
	if rpcClient == nil {
		return nil, errors.New("cannot gen rpcClient of %s", toId)
	}
	res, err := rpcClient.Msg(ctx, req)
	if err != nil {
		return nil, err
	}

	resVal, err := url.ParseQuery(res.GetData())
	return resVal, err
}

func (g *GrpcWorker) getClient(targetWorkerId string) pingablepb.PingableInterfaceClient {
	if client, exist := g.rpcClientMap[targetWorkerId]; exist {
		return client
	}
	if targetWorker, exist := g.worker.ClusterMembers[targetWorkerId]; exist {
		conn, err := grpc.Dial(targetWorker.ServiceAddr, grpc.WithInsecure())
		if err != nil {
			log.Errorf("did not connect: %v", err)
			return nil
		}
		client := pingablepb.NewPingableInterfaceClient(conn)
		g.rpcClientMap[targetWorkerId] = client
		return client
	}
	log.Errorf("cannot gen grpcClient because targetWorker %s not exist", targetWorkerId)
	return nil
}

func (g *GrpcWorker) Serve() error {
	// 开启RPC服务
	lis, err := net.Listen("tcp", g.worker.ServiceAddr)
	if err != nil {
		log.Errorf("监听gRPC端口失败：%v", err)
		return err
	}

	// gRPC通用启动流程
	var opts []grpc.ServerOption
	g.rpcServer = grpc.NewServer(opts...)

	pingablepb.RegisterPingableInterfaceServer(g.rpcServer, &GrpcWorker{})
	reflection.Register(g.rpcServer)

	return g.rpcServer.Serve(lis)
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
			log.Errorf("i(%s) have already follow %s while recv demand follow", w.Id, masterId)
			return nil, nil
		}
		if _, ok := w.ClusterMembers[masterId]; !ok {
			log.Errorf("will follow but masterId:" + masterId + " not exist")
			return nil, nil
		}
		masterPingRes := w.PingNode(masterId)
		if masterPingRes.Code != 0 {
			log.Errorf("will follow(%s) but ping error:" + masterPingRes.Msg)
			return nil, nil
		}
		masterId = masterPingRes.MasterId // follow master's master
		w.Follow(masterId)
		return nil, nil
	})

	return w.pingableWorker.Serve()
}

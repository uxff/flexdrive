package grpcworker

import (
	"context"
	"fmt"
	"net"

	"net/url"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/clusterworker"
	"github.com/uxff/flexdrive/pkg/app/nodestorage/grpcworker/pb/pingablepb"
	"github.com/uxff/flexdrive/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 实现pingableWorker
/**
// 依赖翻转的可操作接口说明
// 消息处理句柄 问题：回到弱类型 但是能兼容grpc和http的实现
type MsgHandler func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error)

// GrpcWorker implements this interface
type PingableWorker interface {
	Serve() error
	PingTo(toId string) (*pingablepb.PingResponse, error) // ping out to other
	//OnPing(MsgHandler)   // like recv ping, grpcServer.Ping instead
	RegisterMsgHandler(action string, handler MsgHandler) // like recv OnMsg
	MsgTo(toId, action, msgId string, param url.Values) (url.Values, error)
	// todo: extend as worker with all functions inlcuding Follow,Remove,Add,EraseMaster,etc
}
**/

// 通信组件 每一个既可以收信又可以发信
// 集成GrpcServer 和 grpc client

type GrpcWorker struct {
	serviceAddr string

	worker *clusterworker.Worker

	rpcServer *grpc.Server

	// need a map of connections
	rpcClientMap map[string]pingablepb.PingableInterfaceClient

	msgHandlerMap map[string]clusterworker.MsgHandler
}

func NewGrpcWorker(worker *clusterworker.Worker) *GrpcWorker {
	return &GrpcWorker{
		worker:        worker,
		rpcClientMap:  make(map[string]pingablepb.PingableInterfaceClient),
		msgHandlerMap: make(map[string]clusterworker.MsgHandler, 0),
	}
}

// on ping, implement grpc
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

// on msg, implement grpc
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

// implement pingableif for clusterworker
func (g *GrpcWorker) PingTo(toId string) (*pingablepb.PingResponse, error) {
	req := &pingablepb.PingRequest{
		FromId:   g.worker.Id,
		MasterId: g.worker.MasterId,
		MetaData: g.worker.WrapMetaData(),
	}

	ctx := context.Background()
	rpcClient := g.getClient(toId)
	if rpcClient == nil {
		return nil, fmt.Errorf("cannot gen rpcClient of %s", toId)
	}
	res, err := rpcClient.Ping(ctx, req)
	return res, err
}

// implement pingableif for clusterworker
func (g *GrpcWorker) RegisterMsgHandler(action string, handler clusterworker.MsgHandler) {
	g.msgHandlerMap[action] = handler
}

// msg out
// implement pingableif for clusterworker
func (g *GrpcWorker) MsgTo(toId, action, msgId string, param url.Values) (url.Values, error) {
	req := &pingablepb.MsgRequest{
		FromId: g.worker.Id,
		ToId:   toId,
		MsgId:  msgId,
		Action: action,
		Data:   param.Encode(),
	}

	ctx := context.Background()
	rpcClient := g.getClient(toId)
	if rpcClient == nil {
		return nil, fmt.Errorf("cannot gen rpcClient of %s", toId)
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
			log.Errorf("can not connect member(%s): %s %v", targetWorkerId, targetWorker.ServiceAddr, err)
			return nil
		}
		client := pingablepb.NewPingableInterfaceClient(conn)
		g.rpcClientMap[targetWorkerId] = client
		return client
	}
	log.Errorf("cannot gen grpcClient because targetWorker %s not exist", targetWorkerId)
	return nil
}

// implement pingableif for clusterworker
func (g *GrpcWorker) Serve(serviceAddr string) error {
	//g.serviceAddr = serviceAddr
	// 开启RPC服务
	lis, err := net.Listen("tcp", serviceAddr)
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

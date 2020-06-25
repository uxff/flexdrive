package grpcpingable

import (
	//"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"net/url"

	"golang.org/x/net/context"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/grpcpingable/pb/pingablepb"
	"github.com/uxff/flexdrive/pkg/app/nodestorage/pingableif"
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
	PingTo(toId string) (url.Values, error) // ping out to other
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

	//worker *clusterworker.Worker

	rpcServer *grpc.Server

	pongHandler pingableif.PongHandler

	// need a map of connections
	rpcClientMap map[string]pingablepb.PingableInterfaceClient
	//rpcClientMap sync.Map

	msgHandlerMap map[string]pingableif.MsgHandler

	lock sync.Mutex
}

func NewGrpcWorker() *GrpcWorker {
	return &GrpcWorker{
		//worker:        worker,
		rpcClientMap:  make(map[string]pingablepb.PingableInterfaceClient),
		msgHandlerMap: make(map[string]pingableif.MsgHandler, 0),
		lock:          sync.Mutex{},
	}
}

// ========== start implement grpc =============
// ping to other, implement grpc
func (g *GrpcWorker) Ping(ctx context.Context, req *pingablepb.PingRequest) (*pingablepb.PingResponse, error) {
	if g.pongHandler != nil {
		resMap, err := g.pongHandler(req.FromId, req.MasterId, req.MetaData)
		res := &pingablepb.PingResponse{
			Code: 0,
			Msg:  "",
			//Members:  w.Members,
			MetaData: resMap.Encode(), // todo rely on interface
		}
		if err != nil {
			//return nil, err
			res.Msg = err.Error()
			res.Code = 1
		}
		return res, nil
	}
	return nil, errors.New("no pong handler registered in GrpcWorker")
}

// on msg, implement grpc Msg(context.Context, *MsgRequest) (*MsgResponse, error)
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

// ========== end implement grpc =============

// ========== start implement pingableif =============

func (g *GrpcWorker) RegisterPong(h pingableif.PongHandler) {
	g.pongHandler = h
}

// implement pingableif for clusterworker
// proto: PingTo(mateAddr string, fromId string, metaData interface{}) (url.Values, error)
func (g *GrpcWorker) PingTo(mateAddr string, fromId string, metaData url.Values) (url.Values, error) {
	req := &pingablepb.PingRequest{
		FromId: fromId,
		//MasterId: g.worker.MasterId,
		MetaData: metaData.Encode(),
	}

	ctx := context.Background()
	rpcClient := g.getClient(mateAddr)
	if rpcClient == nil {
		return nil, fmt.Errorf("cannot gen rpcClient of %s", mateAddr)
	}
	res, err := rpcClient.Ping(ctx, req)
	if err != nil {
		return nil, err
	}

	resVal := url.Values{}
	resVal.Add("masterId", res.MasterId)
	resVal.Add("metaData", res.MetaData)

	return resVal, err
}

// implement pingableif for clusterworker
func (g *GrpcWorker) RegisterMsgHandler(action string, handler pingableif.MsgHandler) {
	g.msgHandlerMap[action] = handler
}

// msg out
// implement pingableif for clusterworker
// proto: MsgTo(mateAddr, action, msgId string, param url.Values) (url.Values, error)
func (g *GrpcWorker) MsgTo(mateAddr, action, msgId string, param url.Values) (url.Values, error) {
	req := &pingablepb.MsgRequest{
		FromId: "", // useless?
		//ToId:   toId,
		MsgId:  msgId,
		Action: action,
		Data:   param.Encode(),
	}

	ctx := context.Background()
	rpcClient := g.getClient(mateAddr)
	if rpcClient == nil {
		return nil, fmt.Errorf("cannot gen rpcClient of %s", mateAddr)
	}
	res, err := rpcClient.Msg(ctx, req)
	if err != nil {
		return nil, err
	}

	resVal, err := url.ParseQuery(res.GetData())
	return resVal, err
}

func (g *GrpcWorker) getClient(targetServiceAddr string) pingablepb.PingableInterfaceClient {
	if client, exist := g.rpcClientMap[targetServiceAddr]; exist {
		return client
	}

	conn, err := grpc.Dial(targetServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("can not connect :%s %v", targetServiceAddr, err)
		return nil
	}
	client := pingablepb.NewPingableInterfaceClient(conn)

	g.lock.Lock()
	defer g.lock.Unlock()

	g.rpcClientMap[targetServiceAddr] = client
	return client
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

// ========== end implement pingableif =============

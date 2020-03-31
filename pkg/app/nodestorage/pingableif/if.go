package pingableif

import (
	"net/url"
)

type MsgHandler func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error)
type PongHandler func(fromId, toId, metaData string) (url.Values, error)

// its a communicate middle ware, clusterworker rely use this interface
// GrpcWorker implements this interface
// HttpWorker implements this interface
type PingableWorker interface {
	Serve(serviceAddr string) error
	PingTo(mateAddr string, fromId string, metaData url.Values) (url.Values, error) // ping out to other
	// like recv ping, cannot use grpcServer.Ping instead
	RegisterPong(PongHandler)
	MsgTo(mateAddr, action, msgId string, param url.Values) (url.Values, error)
	RegisterMsgHandler(action string, handler MsgHandler) // like recv OnMsg
	// todo: extend as worker with all functions inlcuding Follow,Remove,Add,EraseMaster,etc
}

// TODO 内部节点间通信走正式RPC， 外部业务通信可以走伪RPC

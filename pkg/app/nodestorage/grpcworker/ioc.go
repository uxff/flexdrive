package httpworker

// 依赖翻转的可操作接口说明

type PingableHandler interface {
	OnRegistered(w *Worker)
	//MsgTo(jsonableData interface{})
	OnMsg(fromId, data string)
}

// todo ?
type MetaNode interface {
	GetMeta(metaKey string) interface{}
	SetMeta(metaKey string, metaValue interface{})
}

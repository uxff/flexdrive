package httppingable

// tobe instead httpworker

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/app/nodestorage/grpcpingable/pb/pingablepb"
	"github.com/uxff/flexdrive/pkg/app/nodestorage/pingableif"
	"github.com/uxff/flexdrive/pkg/log"
)

// 实现pingableWorker
/**
// 依赖翻转的可操作接口说明
// 消息处理句柄 问题：回到弱类型 但是能兼容grpc和http的实现
type MsgHandler func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error)

// GrpcWorker implements this interface
type PingableWorker interface {
	Serve(serviceAddr string) error
	PingTo(mateAddr string, fromId string, metaData url.Values) (url.Values, error) // ping out to other
	// like recv ping, cannot use grpcServer.Ping instead
	RegisterPong(PongHandler)
	MsgTo(mateAddr, action, msgId string, param url.Values) (url.Values, error)
	RegisterMsgHandler(action string, handler MsgHandler) // like recv OnMsg
	// todo: extend as worker with all functions inlcuding Follow,Remove,Add,EraseMaster,etc
}
**/

type HttpPingableWorker struct {
	serviceAddr string
	router      *gin.Engine

	//worker *clusterworker.Worker

	// need a map of connections
	//rpcClientMap map[string]pingablepb.PingableInterfaceClient

	msgHandlerMap map[string]pingableif.MsgHandler
	pongHandler   pingableif.PongHandler
}

func NewHttpPingableWorker() *HttpPingableWorker {
	return &HttpPingableWorker{
		//rpcClientMap:  make(map[string]pingablepb.PingableInterfaceClient),
		msgHandlerMap: make(map[string]pingableif.MsgHandler, 0),
	}
}

func (w *HttpPingableWorker) PingTo(mateAddr string, fromId string, metaData url.Values) (url.Values, error) {
	res := &pingablepb.PingResponse{}

	req := &pingablepb.PingRequest{
		//MasterId: w.worker.MasterId,
		FromId:   fromId,
		MetaData: metaData.Encode(),
	}

	method := "ping"

	reqBuf, _ := json.Marshal(req)
	reqBufReader := bytes.NewReader(reqBuf) //strings.NewReader(req.Encode()) //

	targetUrl := w.genServeUrl(mateAddr, method)
	resp, err := http.Post(targetUrl, gin.MIMEPOSTForm, reqBufReader)
	if err != nil {
		log.Errorf("ping error:%v", err)
		res.Code, res.Msg = 1, err.Error()
		return nil, err
	}

	resBuf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	//res := url.ParseQuery(string(resBuf))
	err = json.Unmarshal(resBuf, res)
	if err != nil {
		res.Code = 1
		res.Msg = err.Error()
		return nil, err
	}

	resVal := url.Values{}
	resVal.Add("masterId", res.MasterId)
	resVal.Add("metaData", res.MetaData)

	return resVal, nil
}

func (w *HttpPingableWorker) RegisterPong(h pingableif.PongHandler) {
	w.pongHandler = h
}

// proto: MsgTo(mateAddr, action, msgId string, param url.Values) (url.Values, error)
func (w *HttpPingableWorker) MsgTo(mateAddr, action, msgId string, param url.Values) (url.Values, error) {
	res := &pingablepb.MsgResponse{}
	req := &pingablepb.MsgRequest{
		FromId: param.Get("fromId"),
		ToId:   param.Get("toId"),
		Action: action,
		Data:   param.Encode(),
	}
	// req := param
	// req.Add("action", action)
	// req.Add("msgId", "")

	method := "msg"

	reqBuf, _ := json.Marshal(req)
	reqBufReader := bytes.NewReader(reqBuf) //strings.NewReader(req.Encode()) //

	targetUrl := w.genServeUrl(mateAddr, method)
	resp, err := http.Post(targetUrl, gin.MIMEJSON, reqBufReader)
	if err != nil {
		log.Errorf("ping error:%v", err)
		res.Code = 1
		res.Msg = err.Error()
		return nil, err
	}

	resBuf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		res.Code = 1
		res.Msg = err.Error()
		return nil, err
	}

	resVal, err := url.ParseQuery(string(resBuf))

	return resVal, err
}

func (w *HttpPingableWorker) genServeUrl(serviceAddr, method string) string {
	u := url.URL{
		Scheme: "http",
		Host:   serviceAddr,
		Path:   "/" + method,
		//RawQuery: params.Encode(),
	}
	return u.String()
}

func (w *HttpPingableWorker) Serve(serviceAddr string) error {

	router := gin.New() //gin.Default()
	w.router = router

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 接收对方的ping 表示良好 post -d url.Values @expect url.Values
	router.POST("/ping", func(c *gin.Context) {

		res := &pingablepb.PingResponse{}

		b, err := c.GetRawData()
		if err != nil {
			log.Debugf("no error posted")
			res.Code, res.Msg = 1, err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		req := &pingablepb.PingRequest{}
		err = json.Unmarshal(b, req)
		if err != nil {
			log.Debugf("no error posted")
			res.Code, res.Msg = 1, err.Error()
			c.JSON(http.StatusOK, res)
			return
		}

		if w.pongHandler != nil {
			reqVal, _ := url.ParseQuery(req.MetaData)
			resVal, err := w.pongHandler(req.FromId, "", reqVal)
			if err != nil || res == nil {
				log.Debugf("pongHandler error:%v", err)
				res.Code, res.Msg = 1, err.Error()
				c.JSON(http.StatusOK, res)
				return
			}
			res.MasterId = resVal.Get("masterId")
			//res.ToId = resVal.Get("toId")
			res.MetaData = resVal.Encode()
			c.JSON(http.StatusOK, res)
			return
		}
		res.Code, res.Msg = 1, "no pong handler registered in HttpPingaleWorker"

		c.JSON(http.StatusOK, res)
	})

	// 收到消息
	router.POST("/msg", func(c *gin.Context) {
		res := &pingablepb.MsgResponse{
			//MsgId:  req.MsgId,
			//Action: req.Action,
		}
		b, err := c.GetRawData()
		if err != nil {
			log.Debugf("no data posted: %v", err)
			res.Code, res.Msg = 1, err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		req := &pingablepb.MsgRequest{}
		err = json.Unmarshal(b, req)
		if err != nil {
			log.Debugf("illegal data posted:%v", err)
			res.Code, res.Msg = 1, err.Error()
			c.JSON(http.StatusOK, res)
			return
		}

		res.MsgId, res.Action = req.MsgId, req.Action

		if handler, exist := w.msgHandlerMap[req.Action]; exist {
			reqVal, err := url.ParseQuery(req.Data)
			if err != nil {
				res.Code, res.Msg = 1, err.Error()
				c.JSON(http.StatusOK, res)
				return // res, nil
			}
			resVal, err := handler(req.FromId, req.ToId, req.MsgId, reqVal)
			if err != nil {
				res.Code, res.Msg = 1, err.Error()
				c.JSON(http.StatusOK, res)
				return // res, nil
			}
			if resVal != nil {
				res.Data = resVal.Encode()
			}
			c.JSON(http.StatusOK, res)
			return // res, nil
		}

		log.Debugf("when recv msg, no handler regisgered: %v", req)
		res.Code = 99
		res.Msg = "no msg handler registered, ignore"

		c.JSON(http.StatusOK, res)
		//return res, nil
	})

	return router.Run(serviceAddr)
}

func (w *HttpPingableWorker) RegisterMsgHandler(action string, handler pingableif.MsgHandler) {
	log.Debugf("a action registered: action:%s", action)
	w.msgHandlerMap[action] = handler
}

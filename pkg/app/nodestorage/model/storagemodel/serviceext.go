package storagemodel

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	//worker "github.com/uxff/flexdrive/pkg/app/nodestorage/httpworker"
	worker "github.com/uxff/flexdrive/pkg/app/nodestorage/clusterworker"
	"github.com/uxff/flexdrive/pkg/dao/base"

	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/log"
)

const (
	RouteFile = "/file/" // format /file/:fileHash
)

func (n *NodeStorage) AttachService() {

	// 保存文件 依赖注入
	n.Worker.GetPingableWorker().RegisterMsgHandler("savefile", func(fromId, toId, msgId string, reqParam url.Values) (url.Values, error) {

		msg := &NodeMsgSaveFile{}
		//msg.Action = ""
		msg.FileIndexId, _ = strconv.Atoi(reqParam.Get("fileIndexId"))
		msg.FromId = fromId
		msg.AsNodeLevel = reqParam.Get("asNodeLevel")

		err := n.HandleSaveFile(msg)
		if err != nil {
			log.Errorf("handle msg error:%v", err)
			return nil, err
		}

		return nil, nil

	})
}

//
type NodeMsg struct {
	FromId       string
	Action       string
	CustomerAddr string
}

type NodeMsgSaveFile struct {
	NodeMsg
	FileIndexId int
	AsNodeLevel string
}

func (n *NodeStorage) HandleSaveFile(msg *NodeMsgSaveFile) error {

	if msg.FileIndexId == 0 {
		log.Errorf("when handle saveFile, fileIndexId cannot be 0")
		return errors.New("when handle saveFile, fileIndexId cannot be 0")
	}

	fromNode := n.Worker.ClusterMembers[msg.FromId]
	if fromNode == nil {
		//w.JsonError(c, "fromId has no real node")
		return errors.New("fromId has no real node")
	}

	fileIndexEnt, err := dao.GetFileIndexById(msg.FileIndexId)
	if err != nil {
		//w.JsonError(c, "getFileIndex "+fileIndexIdStr+" error")
		log.Errorf("get fileIndexId(%d) error:%v", msg.FileIndexId, err)
		return err
	}

	if fileIndexEnt == nil {
		//w.JsonError(c, "cannot find fileIndexEnt")
		return errors.New("cannot find fileIndexEnt")
	}

	//fileTargetUrl := fromNode.ServiceAddr + "/file/" + fileIndexEnt.FileHash + "/" + fileIndexEnt.FileName
	_, err = n.SaveFileFromFileIndex(msg.FileIndexId, msg.AsNodeLevel)

	return err
}

// useful?
func (n *NodeStorage) OnRegistered(w *worker.Worker) {
	//node.RegisterTo
	n.NodeEnt.LastRegistered = time.Now()
	n.NodeEnt.Status = base.StatusNormal
	n.NodeEnt.UpdateById([]string{"lastRegistered", "status"})
	// todo kick node who timeout
}

// 要求同伴保存文件
func (n *NodeStorage) DemandMateSaveFile(mateId string, fileIndexId int, asNodeLevel string) {
	// msg := &NodeMsgSaveFile{
	// 	NodeMsg: NodeMsg{ // no effected
	// 		FromId: n.Worker.Id,
	// 		Action: "savefile",
	// 	},
	// 	FileIndexId: fileIndexId,
	// }
	urlVal := url.Values{}
	urlVal.Add("fromId", n.Worker.Id)
	urlVal.Add("fileIndexId", strconv.Itoa(fileIndexId))
	urlVal.Add("asNodeLevel", asNodeLevel)
	//val, _ := json.Marshal(msg)
	_, err := n.Worker.GetPingableWorker().MsgTo(n.Worker.ClusterMembers[mateId].ServiceAddr, "savefile", "", urlVal)
	if err != nil {
		log.Errorf("demandMateSaveFile(%s, %d) failed:%v", mateId, fileIndexId, err)
	}
}

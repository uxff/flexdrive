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
		log.Debugf("i(%s) am demanded to savefile: fromId:%s msgId:%s urlVals:%+v", n.Worker.Id, fromId, msgId, reqParam)

		fileIndexId, _ := strconv.Atoi(reqParam.Get("fileIndexId"))
		if fileIndexId == 0 {
			log.Errorf("when handle saveFile, fileIndexId cannot be 0")
			return nil, errors.New("when handle saveFile, fileIndexId cannot be 0")
		}

		fromNode := n.Worker.ClusterMembers[fromId]
		if fromNode == nil {
			//w.JsonError(c, "fromId has no real node")
			log.Errorf("when savefile fromId:%s has no real node", fromId)
			return nil, errors.New("fromId has no real node")
		}

		fileIndexEnt, err := dao.GetFileIndexById(fileIndexId)
		if err != nil {
			//w.JsonError(c, "getFileIndex "+fileIndexIdStr+" error")
			log.Errorf("when savefile get fileIndexId(%d) error:%v", fileIndexId, err)
			return nil, err
		}

		if fileIndexEnt == nil {
			//w.JsonError(c, "cannot find fileIndexEnt")
			log.Errorf("when savefile cannot find fileIndexEnt:%d", fileIndexId)
			return nil, errors.New("cannot find fileIndexEnt")
		}

		//fileTargetUrl := fromNode.ServiceAddr + "/file/" + fileIndexEnt.FileHash + "/" + fileIndexEnt.FileName
		asNodeLevel := reqParam.Get("asNodeLevel")
		_, err = n.SaveFileFromFileIndex(fileIndexId, asNodeLevel)

		if err != nil {
			log.Errorf("when savefile %d error:%v", fileIndexId, err)
			return nil, err
		}
		log.Infof("savefile %d on demand ok", fileIndexId)
		return nil, nil
	})
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
	log.Debugf("will demand savefile: %+v", urlVal)
	_, err := n.Worker.GetPingableWorker().MsgTo(n.Worker.ClusterMembers[mateId].ServiceAddr, "savefile", "", urlVal)
	if err != nil {
		log.Errorf("demandMateSaveFile(%s, %d) failed:%v", mateId, fileIndexId, err)
	}
}

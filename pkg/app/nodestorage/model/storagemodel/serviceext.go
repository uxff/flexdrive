package storagemodel

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/buger/jsonparser"

	worker "github.com/uxff/flexdrive/pkg/app/nodestorage/httpworker"
	//worker "github.com/uxff/flexdrive/pkg/app/nodestorage/clusterworker"
	"github.com/uxff/flexdrive/pkg/dao/base"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/log"
)

const (
	RouteFile = "/file/" // format /file/:fileHash
)

func (n *NodeStorage) AttachService() {
	// 保存文件 但是循环依赖
	// TODO 用onmsg代替
	w := n.Worker
	// 备份文件 fileIndexId fromId nodeLevel
	n.Worker.AttachPostAction("/savefile", func(c *gin.Context) {
		fileIndexIdStr := c.Query("fileIndexId")
		if fileIndexIdStr == "" {
			w.JsonError(c, "fileIndexId must no be empty")
			return
		}

		fileIndexId, _ := strconv.Atoi(fileIndexIdStr)
		if fileIndexId <= 0 {
			w.JsonError(c, "fileIndexId must no be empty")
			return

		}

		fromId := c.Query("fromId")
		if fromId == "" {
			w.JsonError(c, "fromId must no be empty")
			return
		}

		fromNode := w.ClusterMembers[fromId]
		if fromNode == nil {
			w.JsonError(c, "fromId has no real node")
			return
		}

		fileIndexEnt, err := dao.GetFileIndexById(fileIndexId)
		if err != nil {
			w.JsonError(c, "getFileIndex "+fileIndexIdStr+" error")
			log.Errorf("get fileIndexId(%d) error:%v", fileIndexId, err)
			return
		}

		if fileIndexEnt == nil {
			//
			w.JsonError(c, "cannot find fileIndexEnt")
			return
		}

		// 拼出外地址
		//fileUrl := fromNode.ServiceAddr + "://" + fileIndexEnt.OuterPath

		// 下载文件
		// 计算校验hash
		// 保存
		//curNode := storagemodel.GetCurrentNode()
		//localSavePath := curNode.SaveFile()

		//w.masterGoneChan <- true
		log.Debugf("savefile from: %s", fromId)
		w.JsonOk(c)
	})

	// todo 不应该在此处服务
	// n.Worker.AttachGetAction("/file/:fileHash", func(c *gin.Context) {
	// 	fileHash := c.Param("fileHash")
	// 	if fileHash == "" {
	// 		c.String(http.StatusBadRequest, "")
	// 		return
	// 	}

	// 	localFilePath := n.FileHashToStoragePath(fileHash)
	// 	if !DirExist(localFilePath) {
	// 		c.String(http.StatusNotFound, "")
	// 		return
	// 	}

	// 	c.File(localFilePath)
	// })
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

// 依赖注入的设计方式 收到消息的回调
func (n *NodeStorage) OnMsg(fromId, data string) {
	action, err := jsonparser.GetString([]byte(data), "action")
	if err != nil {
		log.Errorf("parse msg error:%v", err)
		return
	}
	switch action {
	case "savefile":
		msg := &NodeMsgSaveFile{}
		err = json.Unmarshal([]byte(data), msg)
		if err != nil {
			log.Errorf("unmarshal msg error:%v", err)
			return
		}
		err = n.HandleSaveFile(msg)
		if err != nil {
			log.Errorf("handle msg error:%v", err)
			return
		}
	default:
		log.Warnf("no action when handle msg:%s", data)
	}
}

func (n *NodeStorage) OnRegistered(w *worker.Worker) {
	//node.RegisterTo
	n.NodeEnt.LastRegistered = time.Now()
	n.NodeEnt.Status = base.StatusNormal
	n.NodeEnt.UpdateById([]string{"lastRegistered", "status"})
	// todo kick node who timeout
}

// 要求同伴保存文件
func (n *NodeStorage) DemandMateSaveFile(mateId string, fileIndexId int, asNodeLevel string) {
	msg := &NodeMsgSaveFile{
		NodeMsg: NodeMsg{ // no effected
			FromId: n.Worker.Id,
			Action: "savefile",
		},
	}
	val, _ := json.Marshal(msg)
	n.Worker.MsgTo(mateId, string(val))
}

package storagemodel

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/log"
)

func (n *NodeStorage) AttachService() {
	// 保存文件 但是循环依赖
	w := n.Worker
	// 备份文件
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

}

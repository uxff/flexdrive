package handler

import (
	"fmt"
	"net/http"

	"github.com/uxff/flexdrive/pkg/dao"

	"github.com/gin-gonic/gin"
)

func Fs(c *gin.Context) {

	fileHash := c.Param("fileHash")
	fileName := c.Param("fileName")
	fileIndex, err := dao.GetFileIndexByFileHash(fileHash)
	if err != nil {
		c.Status(http.StatusNotFound)
		//StdErrMsgResponse(c, ErrInternal, "文件读取错误")
		return
	}

	if fileIndex == nil {
		c.Status(http.StatusNotFound)
		//StdErrMsgResponse(c, ErrInternal, "文件不存在或已被删除")

		return
	}

	if fileName == "" {
		fileName = fileIndex.FileName
	}

	// node := storagemodel.GetCurrentNode()
	// if node == nil {
	// 	StdErrMsgResponse(c, ErrInternal, "文件集群状态异常")
	// 	return
	// }

	//physicalFilePath := node.FileHashToStoragePath(fileHash)

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	//fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	//c.File(physicalFilePath)
	c.File(fileIndex.InnerPath)
}

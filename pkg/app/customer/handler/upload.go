package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func Upload(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ShareSearchRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Share, 0)
	var total int64

	if req.Name != "" {
		total, err = base.ListAndCountByCondition(&dao.Share{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
		if err != nil {
			log.Trace(requestId).Errorf("list failed:%v", err)
			StdErrResponse(c, ErrInternal)
			return
		}
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*ShareItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewShareItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "upload/index.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"total":     total,
		"page":      req.Page,
		"pagesize":  req.PageSize,
		"list":      resItems,
		"reqParam":  req,
		"paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
	})

}

func UploadForm(c *gin.Context) {
	uploadFileKey := "file"
	requestId := c.GetString(CtxKeyRequestId)

	parentDir := c.PostForm("parentDir")
	if parentDir == "" {
		//
		parentDir = "/"
	}
	// 尾巴上补充上/
	if len(parentDir) > 1 && parentDir[len(parentDir)-1] != '/' {
		rparentDir += "/"
	}

	header, err := c.FormFile(uploadFileKey)
	if err != nil {
		//ignore
		log.Trace(requestId).Errorf("form file failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 直接按文件名保存 临时保存 完事要删除的
	dst := header.Filename
	// gin 简单做了封装,拷贝了文件流
	if err := c.SaveUploadedFile(header, dst); err != nil {
		// ignore
		log.Trace(requestId).Errorf("form file failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	defer func() {
		os.Remove(dst)
	}()

	// todo
	// 计算用户空间，是否能够上传
	// 计算hash // todo 客户端计算hash
	// 数据库中保存文件记录
	// 复制备份到多个节点上

	userInfo := getLoginInfo(c)

	if userInfo.UserEnt.UsedSpace+1+header.Size/1000 > userInfo.UserEnt.QuotaSpace {
		StdErrMsgResponse(c, ErrInternal, "您的剩余空间不足，无法上传")
		return
	}

	fileIndex := &dao.FileIndex{
		FileHash: "",
		FileName: header.Filename,
	}

	userFile := &dao.UserFile{
		FileIndexId: fileIndex.Id, // todo
		UserId:      userInfo.UserId,
		FilePath:    parentDir,
		FileName:    header.Filename,
		FileHash:    fileIndex.FileHash,
		IsDir:       0,
		Size:        header.Size,
		Space:       header.Size / 1000,
		Status:      1,
	}

	// 同步到其他节点上

	c.Redirect(http.StatusMovedPermanently, RouteUserFileList)
}

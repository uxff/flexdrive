package handler

import (
	"io"
	"net/http"

	"github.com/uxff/flexdrive/pkg/utils/filehash"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/model/storagemodel"

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
		parentDir += "/"
	}

	header, err := c.FormFile(uploadFileKey)
	if err != nil {
		//ignore
		log.Trace(requestId).Errorf("form file failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 查看节点可用空间够不够
	nodeAvailableSpace := storagemodel.GetCurrentNode().GetFreeSpace()
	if header.Size/1024 >= nodeAvailableSpace {
		log.Trace(requestId).Errorf("no enough available space on such node, upload size:%d available: %d ", header.Size/1024, nodeAvailableSpace)
		StdErrResponse(c, ErrInternal)
		return
	}

	headerFileHandle, err := header.Open()
	if err != nil {
		log.Trace(requestId).Errorf("open upload file error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	defer headerFileHandle.Close() // 必须关闭

	// fileHash := ""
	fileHash, err := filehash.CalcFileSha1(headerFileHandle)
	if err != nil {
		log.Trace(requestId).Errorf("calc filehash if uploaded file failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	curNode := storagemodel.GetCurrentNode()

	// 去集群查找fileHash
	fileIndex, err := dao.GetFileIndexByFileHash(fileHash)
	if err != nil {
		log.Trace(requestId).Errorf("get filehash(%s) error:%v", fileHash, err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 文件不存在 则新建
	if fileIndex == nil {
		headerFileHandle.Seek(0, io.SeekStart)
		log.Trace(requestId).Warnf("filehash(%s) not exist, try build new", fileHash)
		// 保存进集群节点
		fileIndex, err = curNode.SaveFileHandler(headerFileHandle, fileHash, header.Filename, header.Size)
		if err != nil {
			log.Trace(requestId).Errorf("save file(%s) to node(%+v) error:%v", fileHash, curNode.NodeEnt, err)
			StdErrResponse(c, ErrInternal)
			return
		}

		fileHash = fileIndex.FileHash
	}

	//os.MkdirAll("/tmp/flexdirve/", os.ModePerm)

	// 实际应用需要storagemodel来保存
	// 直接按文件名保存 临时保存 完事要删除的
	// dst := "/tmp/flexdrive/" + fileHash //header.Filename
	// gin 简单做了封装,拷贝了文件流
	// if err := c.SaveUploadedFile(header, dst); err != nil {
	// 	// ignore
	// 	log.Trace(requestId).Errorf("form file failed:%v", err)
	// 	StdErrResponse(c, ErrInternal)
	// 	return
	// }

	// defer func() {
	// 	os.Remove(dst)
	// }()

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

	// fileIndex := &dao.FileIndex{
	// 	FileHash: "",
	// 	FileName: header.Filename,
	// }

	userFile := &dao.UserFile{
		FileIndexId: fileIndex.Id,
		UserId:      userInfo.UserId,
		FilePath:    parentDir,
		FileName:    header.Filename,
		FileHash:    fileIndex.FileHash,
		IsDir:       0,
		Size:        header.Size,
		Space:       header.Size/1000 + 1,
		Status:      1,
	}

	userFile.MakePathHash()

	// todo 如果用户目录下已经有该文件名 则不上传
	existUserFileCount, err := base.CountByCondition(&dao.UserFile{}, map[string]interface{}{
		"userId=?":   userInfo.UserId,
		"pathHash=?": userFile.PathHash,
		"fileName=?": userFile.FileName,
		"isDir=?":    userFile.IsDir,
		"status":     base.StatusNormal,
	})

	if err != nil {
		log.Trace(requestId).Errorf("count user(%d) %s%s error:%v", userInfo.UserId, userFile.FilePath, userFile.FileName, err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if existUserFileCount > 0 {
		StdErrMsgResponse(c, ErrInternal, "同名文件已存在，不能覆盖上传:"+userFile.FilePath+userFile.FileName)
		return
	}

	// 保存在用户目录
	_, err = base.Insert(userFile)
	if err != nil {
		log.Trace(requestId).Errorf("insert userFile error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	log.Trace(requestId).Debugf("upload success: %+v", *userFile)

	// 更新用户自己的空间使用记录
	go func(userId int) {
		usedSpace := userFile.SumSpace()
		userInfo.UserEnt.UsedSpace = usedSpace
		err := userInfo.UserEnt.UpdateById([]string{"usedSpace"})
		if err != nil {
			log.Trace(requestId).Debugf("update user(%d) usedSpace failed:%v", userInfo.UserEnt.Id, err)
		}
	}(userInfo.UserId)

	// 同步到其他节点上

	c.Redirect(http.StatusMovedPermanently, RouteUserFileList+"?dir="+parentDir)
}

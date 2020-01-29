package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func init() {
}

type FileIndexListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"fileName"`
	FileHash    string `form:"fileHash"`
	NodeId      int    `form:"nodeId"`
	Status      int    `form:"status"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *FileIndexListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.CreateStart != "" {
		condition["created>=?"] = r.CreateStart
	}

	if r.CreateEnd != "" {
		condition["created<=?"] = r.CreateEnd
	}

	if r.Name != "" {
		condition["fileName like ?"] = "%" + r.Name + "%"
	}

	if r.FileHash != "" {
		condition["fileHash = ?"] = r.FileHash
	}

	if r.NodeId != 0 {
		condition["node = ?"] = r.NodeId
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type FileIndexItem struct {
	dao.FileIndex
}

func NewFileIndexItemFromEnt(fileIndexEnt *dao.FileIndex) *FileIndexItem {
	return &FileIndexItem{
		FileIndex: *fileIndexEnt,
	}
}

func FileIndexList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &FileIndexListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.FileIndex, 0)
	total, err := base.ListAndCountByCondition(&dao.FileIndex{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*FileIndexItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewFileIndexItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "fileindex/list.tpl", gin.H{
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

func FileIndexEnable(c *gin.Context) {
	fileIndexId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if fileIndexId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	//loginInfo := getLoginInfo(c)

	shareEnt, err := dao.GetFileIndexById(int(fileIndexId))

	//_, err := base.GetByCol("id", mid, shareEnt)
	// exist, err := base.GetByCol("mid", mid, shareEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if shareEnt == nil {
		StdErrResponse(c, ErrUserNotExist)
		return
	}

	if enable == 1 {
		// 启用
		shareEnt.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		shareEnt.Status = base.StatusDeleted
	}

	_, err = base.UpdateByCol("id", fileIndexId, shareEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteFileIndexList)
}

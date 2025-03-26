package apihandler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

func init() {
}

type FileIndexListRequest struct {
	CreateStart string `json:"createStart"`
	CreateEnd   string `json:"createEnd"`
	Name        string `json:"fileName"`
	FileHash    string `json:"fileHash"`
	NodeId      int    `json:"nodeId"`
	Status      int    `json:"status"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pagesize"`
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
		condition["nodeId = ?"] = r.NodeId
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

// FileIndexList - to list the FileIndex
func FileIndexList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &FileIndexListRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		JsonErr(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.FileIndex, 0)
	total, err := base.ListAndCountByCondition(&dao.FileIndex{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*FileIndexItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewFileIndexItemFromEnt(v))
	}

	JsonOk(c, gin.H{
		"total":    total,
		"page":     req.Page,
		"pagesize": req.PageSize,
		"list":     resItems,
		"reqParam": req,
		// "paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

// FileIndexEnable - to enable the FileIndex
func FileIndexEnable(c *gin.Context) {
	fileIndexId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if fileIndexId <= 0 {
		JsonErr(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	//loginInfo := getLoginInfo(c)

	fileIndexEnt, err := dao.GetFileIndexById(int(fileIndexId))

	//_, err := base.GetByCol("id", mid, fileIndexEnt)
	// exist, err := base.GetByCol("mid", mid, fileIndexEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	if fileIndexEnt == nil {
		JsonErr(c, ErrMgrNotExist)
		return
	}

	if enable == 1 {
		// 启用
		fileIndexEnt.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		fileIndexEnt.Status = base.StatusDeleted
	}

	_, err = base.UpdateByCol("id", fileIndexId, fileIndexEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	JsonOk(c, gin.H{
		"id": fileIndexEnt.Id,
	})
	// c.Redirect(http.StatusMovedPermanently, RouteFileIndexList)
}

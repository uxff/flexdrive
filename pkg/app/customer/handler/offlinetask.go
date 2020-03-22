package handler

import (
	"net/http"
	"strconv"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/model/storagemodel"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func init() {
}

type OfflineTaskListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"name"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *OfflineTaskListRequest) ToCondition() (condition map[string]interface{}) {
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

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type OfflineTaskItem struct {
	dao.OfflineTask
}

func NewOfflineTaskItemFromEnt(offlineTaskEnt *dao.OfflineTask) *OfflineTaskItem {
	return &OfflineTaskItem{
		OfflineTask: *offlineTaskEnt,
	}
}

func OfflineTaskList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &OfflineTaskListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	loginInfo := getLoginInfo(c)
	condition := req.ToCondition()
	condition["userId=?"] = loginInfo.UserId

	// 列表查询
	list := make([]*dao.OfflineTask, 0)
	total, err := base.ListAndCountByCondition(&dao.OfflineTask{}, condition, req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*OfflineTaskItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewOfflineTaskItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "offlinetask/list.tpl", gin.H{
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

type OfflineTaskAddRequest struct {
	ParentUserFileId int    `form:"parentUserFileId"` // 文件
	Dataurl          string `form:"dataurl"`          // 资源地址
	// ExpiredTime time.Time `form:"-"`
}

func OfflineTaskAdd(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &OfflineTaskAddRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		log.Trace(requestId).Errorf("bind param error:%v", err)
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	if req.Dataurl == "" {
		StdErrMsgResponse(c, ErrInvalidParam, "文件id为空")
		return
	}

	loginInfo := getLoginInfo(c)

	offlineTaskItem := &dao.OfflineTask{
		UserId:           loginInfo.UserId,
		Dataurl:          req.Dataurl,
		ParentUserFileId: req.ParentUserFileId,
		//UserFileId: req.UserFileId,
		//FileName:   userFile.FileName,
		//FileHash:   userFile.FileHash,
		//NodeId:     userFile.NodeId,
		//FilePath:   userFile.FilePath,
		Status: base.StatusNormal,
		//Expired: time.Now().Add(time.Hour * 24 * 7),
	}

	_, err = base.Insert(offlineTaskItem)
	if err != nil {
		log.Trace(requestId).Errorf("db insert error:%v", err)
		StdErrMsgResponse(c, ErrInternal, "创建离线任务失败")
		return
	}

	go func() {
		node := storagemodel.GetCurrentNode()
		if node != nil {
			//startOfflineTask(offlineTaskItem)
			err := node.ExecOfflineTask(offlineTaskItem)
			if err != nil {
				log.Trace(requestId).Errorf("exec offlinetask(%d) failed:%v", offlineTaskItem.Id, err)
			}
		}
	}()

	StdResponse(c, ErrSuccess, offlineTaskItem) // todo redirect
	return

}

// todo 还允许重新开始吗
func OfflineTaskEnable(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)
	offlineTaskId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if offlineTaskId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	offlineTaskEnt, err := dao.GetOfflineTaskById(int(offlineTaskId))

	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if offlineTaskEnt == nil {
		StdErrResponse(c, ErrItemNotExist)
		return
	}

	if enable == 1 {
		// 启用
		offlineTaskEnt.Status = dao.OfflineTaskStatusExecuting

		go func() {
			node := storagemodel.GetCurrentNode()
			if node != nil {
				//startOfflineTask(offlineTaskItem)
				err := node.ExecOfflineTask(offlineTaskEnt)
				if err != nil {
					log.Trace(requestId).Errorf("exec offlinetask(%d) failed:%v", offlineTaskEnt.Id, err)
				}
			}
		}()
	} else if enable == 9 {
		// 停用
		offlineTaskEnt.Status = base.StatusDeleted
	}

	//base.CacheDelByEntity("mgrLoginName", offlineTaskEnt.Email, offlineTaskEnt)

	_, err = base.UpdateByCol("id", offlineTaskId, offlineTaskEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteUserFileList)
}

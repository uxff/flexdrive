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

type UserFileListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"fileName"`
	FileHash    string `form:"fileHash"`
	UserId      int    `form:"userId"`
	Status      int    `form:"status"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *UserFileListRequest) ToCondition() (condition map[string]interface{}) {
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

	if r.UserId != 0 {
		condition["userId = ?"] = r.UserId
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type UserFileItem struct {
	dao.UserFile
}

func NewUserFileItemFromEnt(fileIndexEnt *dao.UserFile) *UserFileItem {
	return &UserFileItem{
		UserFile: *fileIndexEnt,
	}
}

func UserFileList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &UserFileListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.UserFile, 0)
	total, err := base.ListAndCountByCondition(&dao.UserFile{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*UserFileItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewUserFileItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "user/filelist.tpl", gin.H{
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

func UserFileEnable(c *gin.Context) {
	userFileId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if userFileId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	//loginInfo := getLoginInfo(c)

	userFileEnt, err := dao.GetUserFileById(int(userFileId))

	//_, err := base.GetByCol("id", mid, userFileEnt)
	// exist, err := base.GetByCol("mid", mid, userFileEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if userFileEnt == nil {
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if enable == 1 {
		// 启用
		userFileEnt.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		userFileEnt.Status = base.StatusDeleted
	}

	_, err = base.UpdateByCol("id", userFileId, userFileEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteUserFileList+"?userId="+strconv.Itoa(userFileEnt.UserId))
}

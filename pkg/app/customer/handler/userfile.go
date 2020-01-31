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
	FileName    string `form:"fileName"`
	Dir         string `form:"dir"`
	SearchDir   int    `form:"searchDir"` // 是否在当前目录下搜索 默认搜索全部目录
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

	if r.FileName != "" {
		condition["fileName like ?"] = "%" + r.FileName + "%"
	}

	if r.SearchDir == 1 {
		fileIdxTmp := &dao.UserFile{
			FilePath: r.Dir,
		}
		condition["pathHash= ?"] = fileIdxTmp.MakePathHash()
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

	condition := req.ToCondition()
	condition["status=?"] = base.StatusNormal // 只查询未删除

	// 列表查询
	list := make([]*dao.UserFile, 0)
	total, err := base.ListAndCountByCondition(&dao.UserFile{}, condition, req.Page, req.PageSize, "", &list)
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

	c.HTML(http.StatusOK, "userfile/list.tpl", gin.H{
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
	fileIndexId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if fileIndexId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	//loginInfo := getLoginInfo(c)

	shareEnt, err := dao.GetUserFileById(int(fileIndexId))

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
	c.Redirect(http.StatusMovedPermanently, RouteUserFileList)
}

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

type UserListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Email       string `form:"email"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *UserListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.CreateStart != "" {
		condition["created>=?"] = r.CreateStart
	}

	if r.CreateEnd != "" {
		condition["created<=?"] = r.CreateEnd
	}

	// if r.Name != "" {
	// 	condition["name like ?"] = "%" + r.Email + "%"
	// }

	if r.Email != "" {
		condition["email = ?"] = r.Email
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type UserItem struct {
	dao.User
}

func NewUserItemFromEnt(userEnt *dao.User) *UserItem {
	return &UserItem{
		User: *userEnt,
	}
}

func UserList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &UserListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.User, 0)
	total, err := base.ListAndCountByCondition(&dao.User{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*UserItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewUserItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "user/list.tpl", gin.H{
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

func UserEnable(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if userId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	//loginInfo := getLoginInfo(c)

	userEnt, err := dao.GetUserById(int(userId))

	//_, err := base.GetByCol("id", mid, userEnt)
	// exist, err := base.GetByCol("mid", mid, userEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if userEnt == nil {
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if enable == 1 {
		// 启用
		userEnt.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		userEnt.Status = base.StatusDeleted
	}

	//base.CacheDelByEntity("mgrLoginName", userEnt.Email, userEnt)

	_, err = base.UpdateByCol("id", userId, userEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteUserList)
}

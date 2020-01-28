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
	Id          int    `json:"id" form:"id"`
	Email       string `json:"email" form:"email" binding:"required"`
	LastLoginAt string `json:"lastLoginAt"`
	LastLoginIp string `json:"lastLoginIp"`
	//Pwd         string `json:"pwd" form:"pwd"`
	LevelId     int    `json:"levelId"`
	LevelName   string `json:"levelName"`
	TotalCharge int    `json:"totalCharge"`
	QuotaSpace  int64  `json:"quotaSpace"`
	UsedSpace   int64  `json:"usedSpace"`
	FileCount   int64  `json:"fileCount"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	Status      int    `json:"status"`
}

func NewUserItemFromEnt(userEnt *dao.User) *UserItem {
	return &UserItem{
		Id:          userEnt.Id,
		Email:       userEnt.Email,
		LastLoginAt: userEnt.LastLoginAt.String(),
		LastLoginIp: userEnt.LastLoginIp,
		//RoleId:      mgruserEntEnt.RoleId,
		// Pwd:         userEnt.Pwd, // 不返回密码，如果更新时提交了密码则代表修改密码
		//RoleName: userEnt.RoleName,
		FileCount:  userEnt.FileCount,
		QuotaSpace: userEnt.QuotaSpace,
		UsedSpace:  userEnt.UsedSpace,
		Created:    userEnt.Created.String(),
		Updated:    userEnt.Updated.String(),
		Status:     userEnt.Status,
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

	// StdResponse(c, ErrSuccess, map[string]interface{}{
	// 	"total":    total,
	// 	"page":     req.Page,
	// 	"pagesize": req.PageSize,
	// 	"data":     resItems,
	// })
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

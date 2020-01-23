package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

type ManagerListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Email       string `form:"email"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *ManagerListRequest) ToCondition() (condition map[string]interface{}) {
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
type ManagerItem struct {
	Mid         int    `json:"mid"`
	Email       string `json:"email"`
	LastLoginAt string `json:"lastLoginAt"`
	LastLoginIp string `json:"lastLoginIp"`
	Pwd         string `json:"pwd"`
	RoleId      int    `json:"roleId"`
	RoleName    string `json:"roleName"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	Status      int    `json:"status"`
}

func NewManagerItemFromEnt(mgrEnt *dao.Manager) *ManagerItem {
	return &ManagerItem{
		Mid:         mgrEnt.Id,
		Email:       mgrEnt.Email,
		LastLoginAt: mgrEnt.LastLoginAt.String(),
		LastLoginIp: mgrEnt.LastLoginIp,
		RoleId:      mgrEnt.RoleId,
		// Pwd:         mgrEnt.Pwd, // 不返回密码，如果更新时提交了密码则代表修改密码
		//RoleName: mgrEnt.RoleName,
		Created: mgrEnt.Created.String(),
		Updated: mgrEnt.Updated.String(),
		Status:  mgrEnt.Status,
	}
}

func ManagerList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ManagerListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Manager, 0)
	total, err := base.ListAndCountByCondition(&dao.Manager{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*ManagerItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewManagerItemFromEnt(v))
	}

	// StdResponse(c, ErrSuccess, map[string]interface{}{
	// 	"total":    total,
	// 	"page":     req.Page,
	// 	"pagesize": req.PageSize,
	// 	"data":     resItems,
	// })
	c.HTML(http.StatusOK, "manager/list.tpl", gin.H{
		"total":    total,
		"page":     req.Page,
		"pagesize": req.PageSize,
		"data":     resItems,
	})
}

type ManagerAddRequest struct {
	ManagerItem // 新增只用到里面的 Email,roleid,pwd
	// Status int
}

func (r *ManagerAddRequest) ToEnt() *dao.Manager {
	e := &dao.Manager{
		Email:  r.Email,
		RoleId: r.RoleId,
		// MgrLastLoginAt:time.Now(),
		//Pwd: r.Pwd,
	}
	e.SetPwd(r.Pwd)
	return e
}

func ManagerAdd(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	req := &ManagerAddRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 去掉用户输入的字符串里开头结尾的不可见字符
	req.Email = strings.TrimSpace(req.Email)
	req.Pwd = strings.TrimSpace(req.Pwd)
	req.RoleName = strings.TrimSpace(req.RoleName)

	log.Trace(requestId).Debugf("%+v", req)

	// 检验名称是否已经存在
	existEnt, err := dao.GetManagerByEmail(req.Email)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if existEnt != nil && existEnt.Id != req.Mid {
		StdErrResponse(c, ErrNameDuplicate)
		return
	}

	mid := req.Mid // 如果有则是编辑
	// roleEnt, err := dao.GetRoleById(req.RoleId)
	// if err != nil {
	// 	log.Trace(requestId).Errorf("查询角色信息失败:%v roleId:%d", err, req.RoleId)
	// 	return
	// }

	mgrEnt := req.ToEnt()
	mgrEnt.LastLoginAt.UnmarshalJSON([]byte(time.Now().String()))
	mgrEnt.LastLoginIp = c.Request.Header.Get("X-Real-IP")
	mgrEnt.Status = base.StatusNormal
	//mgrEnt.RoleName = roleEnt.Name

	if mgrEnt.Name == "" { // 用户名不能为空
		StdResponseJson(c, ErrInvalidParam, "用户名不能为空", "")
		log.Trace(requestId).Warnf("用户名不能为空")
		return
	}

	if mid > 0 {
		cols := []string{"email", "roleId", "roleName"}
		if req.Pwd != "" { // 默认如果密码不为空，则更新密码
			cols = append(cols, "mgrPwd")
			mgrEnt.SetPwd(req.Pwd)
		}
		_, err = base.UpdateByCol("mid", mid, mgrEnt, cols)
		//base.CacheDelByEntity("mgrLoginName", req.Email, existEnt)
	} else {
		if req.Pwd == "" {
			StdResponseJson(c, ErrInvalidParam, "密码不能为空", "")
			log.Trace(requestId).Warnf("密码不能为空")
			return
		}

		//mgrEnt.Pwd, _ = utils.Md5Sum([]byte(mgrEnt.Pwd))
		_, err = base.Insert(mgrEnt)
		mid = mgrEnt.Id
	}

	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// StdResponse(c, ErrSuccess, gin.H{
	// 	"mid": mid,
	// })
	c.Redirect(http.StatusMovedPermanently, RouteManagerList)

}

func ManagerEnable(c *gin.Context) {
	mid, _ := strconv.ParseInt(c.Param("mid"), 10, 64)
	if mid <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	loginInfo := getLoginInfo(c)

	mgrEnt, err := dao.GetManagerById(int(mid)) //&dao.Manager{}

	//_, err := base.GetByCol("id", mid, mgrEnt)
	// exist, err := base.GetByCol("mid", mid, mgrEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if mgrEnt == nil {
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if enable == 1 {
		// 启用
		mgrEnt.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		if int(mid) == loginInfo.Mid {
			StdResponseJson(c, ErrInternal, "不能停用自己的账号", "")
			return
		}
		if mgrEnt.IsSuperRole() {
			StdResponseJson(c, ErrInternal, "不能停用超级管理的账号", "")
			return
		}
		mgrEnt.Status = base.StatusDeleted
	}

	//base.CacheDelByEntity("mgrLoginName", mgrEnt.Email, mgrEnt)

	_, err = base.UpdateByCol("mid", mid, mgrEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteManagerList)
}

type ManagerChangePwdRequest struct {
	// Email string `json:"mgrName" binding:"required"`
	Oldpwd string `json:"oldpwd" binding:"required"`
	Newpwd string `json:"pwd" binding:"required"`
	// Force int ?
	// Captcha ?
}

// 超管修改别人的手机号
func ManagerChangePwd(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	req := &ManagerChangePwdRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	loginEnt := getLoginInfo(c)
	if loginEnt == nil || loginEnt.Mid <= 0 {
		StdErrResponse(c, ErrNotLogin)
		return
	}

	mgrEnt, err := dao.GetManagerById(loginEnt.Mid)
	//_, err = base.GetByCol("id", loginEnt.Mid, mgrEnt)
	if err != nil {
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if mgrEnt.IsPwdValid(req.Oldpwd) {
		StdErrResponse(c, ErrInvalidPass)
		return
	}

	mgrEnt.SetPwd(req.Newpwd)

	// _, err = base.UpdateByCol("mid", loginEnt.Mid, mgrDbEnt, []string{"mgrPwd"})
	err = mgrEnt.UpdateById([]string{"pwd"})

	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteManagerList)
}

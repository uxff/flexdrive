package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

type RoleListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"name"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *RoleListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.CreateStart != "" {
		condition["created>=?"] = r.CreateStart
	}

	if r.CreateEnd != "" {
		condition["created<=?"] = r.CreateEnd
	}

	if r.Name != "" {
		condition["name like ?"] = "%" + r.Name + "%"
	}

	// if r.Email != "" {
	// 	condition["email = ?"] = r.Email
	// }

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type RoleItem struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Permit  string `json:"permit"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Status  int    `json:"status"`
}

func RoleList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ManagerListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Role, 0)
	total, err := base.ListAndCountByCondition(&dao.Role{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	//resItems := make([]*Role, 0)
	// for _, roleItem := range list {
	//resItems = append(resItems, NewManagerItemFromEnt(v))
	// parse permit
	// }

	StdResponse(c, ErrSuccess, map[string]interface{}{
		"total":    total,
		"page":     req.Page,
		"pagesize": req.PageSize,
		"data":     list,
	})
}

type RoleAddRequest struct {
	RoleItem // 新增只用到里面的 Email,roleid,pwd
	// Status int
}

func (r *RoleAddRequest) ToEnt() *dao.Role {
	e := &dao.Role{
		Name: r.Name,
		//RoleId: r.RoleId,
		// MgrLastLoginAt:time.Now(),
		//Pwd: r.Pwd,
	}
	//e.SetPwd(r.Pwd)
	return e
}

// 新增和修改
func RoleAdd(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	req := &RoleAddRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	if req.Name == "" { // 用户名不能为空
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 去掉用户输入的字符串里开头结尾的不可见字符
	req.Name = strings.TrimSpace(req.Name)
	req.Permit = strings.TrimSpace(req.Permit)

	log.Trace(requestId).Debugf("%+v", req)

	// 检验名称是否已经存在
	existEnt, err := dao.GetRoleByName(req.Name)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 如果是添加 不能添加同名
	if existEnt != nil && req.Id == 0 {
		StdErrResponse(c, ErrNameDuplicate)
		return
	}

	// 如果是修改 不能修改和已存在的冲突
	if existEnt != nil && req.Id > 0 && req.Id != existEnt.Id {
		StdErrResponse(c, ErrNameDuplicate)
		return
	}

	// mid := req.Mid // 如果有则是编辑
	// roleEnt, err := dao.GetRoleById(req.RoleId)
	// if err != nil {
	// 	log.Trace(requestId).Errorf("查询角色信息失败:%v roleId:%d", err, req.RoleId)
	// 	return
	// }

	ent := req.ToEnt()
	ent.Status = base.StatusNormal

	if req.Id > 0 {
		cols := []string{"name"}
		_, err = base.UpdateByCol("id", req.Id, ent, cols)
		//base.CacheDelByEntity("mgrLoginName", req.Email, existEnt)
	} else {
		_, err = base.Insert(ent)
	}

	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// StdResponse(c, ErrSuccess, gin.H{
	// 	"mid": mid,
	// })
	c.Redirect(http.StatusMovedPermanently, RouteRoleList)

}

func RoleEnable(c *gin.Context) {
	roleId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if roleId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	loginInfo := getLoginInfo(c)

	ent, err := dao.GetRoleById(int(roleId))

	//_, err := base.GetByCol("id", mid, mgrEnt)
	// exist, err := base.GetByCol("mid", mid, mgrEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if ent == nil {
		StdErrResponse(c, ErrRoleNotExist)
		return
	}

	if enable == 1 {
		// 启用
		ent.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		if int(roleId) == loginInfo.RoleId {
			StdResponseJson(c, ErrInternal, "不能停用自己的账号", "")
			return
		}
		// if ent.IsSuperRole() {
		// 	StdResponseJson(c, ErrInternal, "不能停用超级管理的账号", "")
		// 	return
		// }
		ent.Status = base.StatusDeleted
	}

	//base.CacheDelByEntity("mgrLoginName", mgrEnt.Email, mgrEnt)

	_, err = base.UpdateByCol("id", roleId, ent, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteRoleList)
}

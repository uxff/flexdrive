package apihandler

import (
	"strconv"
	"strings"

	"github.com/uxff/flexdrive/pkg/app/admin/model/rbac"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
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

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type RoleItem struct {
	Id          int               `json:"id"`
	Name        string            `json:"name" form:"name"`
	Created     string            `json:"created" form:"-"`
	Updated     string            `json:"updated" form:"-"`
	Status      int               `json:"status" form:"-"`
	AccessRoute map[string]string `json:"accessRoute" form:"accessRoute"`
}

func RoleList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &RoleListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		JsonErr(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Role, 0)
	total, err := base.ListAndCountByCondition(&dao.Role{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	//resItems := make([]*Role, 0)
	// for _, roleItem := range list {
	//resItems = append(resItems, NewManagerItemFromEnt(v))
	// parse permit
	// }

	JsonOk(c, gin.H{
		"total":     total,
		"page":      req.Page,
		"pagesize":  req.PageSize,
		"list":      list,
		"reqParam":  req,
		"paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
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

	e.Permit = rbac.GetAllMenu()

	for _, accessGroup := range e.Permit {
		// 一级分组不用配置权限
		for _, accessItem := range accessGroup.Sub {
			accessItem.Access = false
			for routeStr, isAccess := range r.AccessRoute {
				if routeStr == accessItem.PermitRoute && isAccess == "1" {
					accessItem.Access = true
				}
			}
		}
	}

	//e.SetPwd(r.Pwd)
	return e
}

// func RoleAdd(c *gin.Context) {
// 	allAccessItems := rbac.GetAllMenu()
// 	c.HTML(http.StatusOK, "role/add.tpl", gin.H{
// 		"LoginInfo":      getLoginInfo(c),
// 		"IsLogin":        isLoginIn(c),
// 		"allAccessItems": allAccessItems,
// 	})
// }

// func RoleEdit(c *gin.Context) {
// 	roleId, _ := strconv.Atoi(c.Param("id"))
// 	if roleId <= 0 {
// 		JsonErr(c, ErrInvalidParam)
// 		return
// 	}

// 	roleEnt, _ := dao.GetRoleById(roleId)
// 	if roleEnt == nil {
// 		JsonErr(c, ErrMgrNotExist)
// 		return
// 	}

// 	if roleEnt.Status == base.StatusDeleted {
// 		JsonErr(c, ErrMgrDisabled)
// 		return
// 	}

// 	log.Debugf("roleEnt:%+v", roleEnt)
// 	allAccessItems := rbac.GetAllAccessItems(roleEnt.Permit)

// 	c.HTML(http.StatusOK, "role/edit.tpl", gin.H{
// 		"LoginInfo":      getLoginInfo(c),
// 		"IsLogin":        isLoginIn(c),
// 		"roleId":         roleId,
// 		"RoleEnt":        roleEnt,
// 		"allAccessItems": allAccessItems,
// 	})
// }

// 新增和修改
func RoleAddForm(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)
	roleId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := &RoleAddRequest{}
	err := c.ShouldBind(req)
	req.AccessRoute = c.PostFormMap("accessRoute")
	if err != nil {
		JsonErr(c, ErrInvalidParam)
		return
	}

	if req.Name == "" { // 角色名不能为空
		JsonErr(c, ErrInvalidParam)
		return
	}

	// 去掉用户输入的字符串里开头结尾的不可见字符
	req.Name = strings.TrimSpace(req.Name)
	//req.Permit = strings.TrimSpace(req.Permit)

	log.Trace(requestId).Debugf("%+v", req)

	// 检验名称是否已经存在
	existEnt, err := dao.GetRoleByName(req.Name)
	if err != nil {
		log.Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	// 如果是添加 不能添加同名
	if existEnt != nil && roleId == 0 {
		JsonErr(c, ErrNameDuplicate)
		return
	}

	// 如果是修改 不能修改和已存在的冲突
	if existEnt != nil && roleId > 0 && int(roleId) != existEnt.Id {
		JsonErr(c, ErrNameDuplicate)
		return
	}

	// roleId := req.roleId // 如果有则是编辑
	// roleEnt, err := dao.GetRoleById(req.RoleId)
	// if err != nil {
	// 	log.Trace(requestId).Errorf("查询角色信息失败:%v roleId:%d", err, req.RoleId)
	// 	return
	// }

	ent := req.ToEnt()
	ent.Status = base.StatusNormal

	if roleId > 0 {
		cols := []string{"name", "permit"}
		_, err = base.UpdateByCol("id", roleId, ent, cols)
		//base.CacheDelByEntity("mgrLoginName", req.Email, existEnt)
	} else {
		_, err = base.Insert(ent)
	}

	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	JsonOk(c, gin.H{
		"id": ent.Id,
	})
	// c.Redirect(http.StatusMovedPermanently, RouteRoleList)

}

func RoleEnable(c *gin.Context) {
	roleId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if roleId <= 0 {
		JsonErr(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	loginInfo := getLoginInfo(c)

	ent, err := dao.GetRoleById(int(roleId))

	//_, err := base.GetByCol("id", roleId, roleEnt)
	// exist, err := base.GetByCol("roleId", roleId, roleEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	if ent == nil {
		JsonErr(c, ErrRoleNotExist)
		return
	}

	if enable == 1 {
		// 启用
		ent.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		if int(roleId) == loginInfo.RoleId {
			JsonErrMsg(c, ErrInternal, "不能停用自己的账号")
			return
		}
		// if ent.IsSuperRole() {
		// 	StdErrMsgResponse(c, ErrInternal, "不能停用超级管理的账号", "")
		// 	return
		// }
		ent.Status = base.StatusDeleted
	}

	_, err = base.UpdateByCol("id", roleId, ent, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	JsonOk(c, gin.H{
		"id": ent.Id,
	})
	// c.Redirect(http.StatusMovedPermanently, RouteRoleList)
}

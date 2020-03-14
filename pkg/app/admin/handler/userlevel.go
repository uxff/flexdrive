package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

type UserLevelListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"name"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *UserLevelListRequest) ToCondition() (condition map[string]interface{}) {
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
type UserLevelItem struct {
	Id         int    `json:"id"`
	Name       string `json:"name" form:"name"`
	QuotaSpace int    `json:"quotaSpace" form:"quotaSpace"` // kB
	Price      int    `json:"price" form:"price"`           // 分
	IsDefault  int    `json:"isDefault" form:"isDefault"`   // 分
	Created    string `json:"created" form:"-"`
	Updated    string `json:"updated" form:"-"`
	Status     int    `json:"status" form:"-"`
	Desc       string `json:"desc" form:"desc"`
	PrimeCost  int    `json:"primeCost" form:"primeCost"` // 分
}

func UserLevelList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &UserLevelListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	req.PageSize = 1000

	// 列表查询
	list := make([]*dao.UserLevel, 0)
	total, err := base.ListAndCountByCondition(&dao.UserLevel{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	//resItems := make([]*UserLevel, 0)
	// for _, UserLevelItem := range list {
	//resItems = append(resItems, NewManagerItemFromEnt(v))
	// parse permit
	// }

	c.HTML(http.StatusOK, "userlevel/list.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"total":     total,
		"page":      req.Page,
		"pagesize":  req.PageSize,
		"list":      list,
		"reqParam":  req,
		"paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

type UserLevelAddRequest struct {
	UserLevelItem // 新增只用到里面的 Email,levelId,pwd
	// Status int
}

func (r *UserLevelAddRequest) ToEnt() *dao.UserLevel {
	e := &dao.UserLevel{
		Name:       r.Name,
		QuotaSpace: int64(r.QuotaSpace),
		Price:      r.Price,
		IsDefault:  r.IsDefault,
		PrimeCost:  r.PrimeCost,
		Desc:       r.Desc,
		//levelId: r.levelId,
		// MgrLastLoginAt:time.Now(),
		//Pwd: r.Pwd,
	}

	//e.SetPwd(r.Pwd)
	return e
}

func UserLevelAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "userlevel/add.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
	})
}

func UserLevelEdit(c *gin.Context) {
	levelId, _ := strconv.Atoi(c.Param("id"))
	if levelId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	levelEnt, _ := dao.GetUserLevelById(levelId)
	if levelEnt == nil {
		StdErrResponse(c, ErrLevelNotExist)
		return
	}

	if levelEnt.Status == base.StatusDeleted {
		StdErrResponse(c, ErrLevelDisabled)
		return
	}

	log.Debugf("levelEnt:%+v", levelEnt)

	c.HTML(http.StatusOK, "userlevel/edit.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"levelId":   levelId,
		"levelEnt":  levelEnt,
	})
}

// 新增和修改
func UserLevelAddForm(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)
	levelId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := &UserLevelAddRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	if req.Name == "" { // 名不能为空
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 去掉用户输入的字符串里开头结尾的不可见字符
	req.Name = strings.TrimSpace(req.Name)
	//req.Permit = strings.TrimSpace(req.Permit)

	log.Trace(requestId).Debugf("%+v", req)

	// 检验名称是否已经存在
	existEnt, err := dao.GetUserLevelByName(req.Name)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 如果是添加 不能添加同名
	if existEnt != nil && levelId == 0 {
		StdErrResponse(c, ErrNameDuplicate)
		return
	}

	// 如果是修改 不能修改和已存在的冲突 似乎有BUG
	if existEnt != nil && levelId > 0 && int(levelId) != existEnt.Id {
		StdErrResponse(c, ErrNameDuplicate)
		return
	}

	ent := req.ToEnt()
	ent.Status = base.StatusNormal

	if levelId > 0 {
		if ent.IsDefault == 1 {
			// 将其他的都设置为非默认
			base.UpdateByCondition(&dao.UserLevel{IsDefault: 0}, map[string]interface{}{"status=?": 1}, []string{"isDefault"})
		}

		cols := []string{"name", "quotaSpace", "price", "isDefault"}
		_, err = base.UpdateByCol("id", levelId, ent, cols)
	} else {
		_, err = base.Insert(ent)
	}

	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	c.Redirect(http.StatusMovedPermanently, RouteUserLevelList)

}

func UserLevelEnable(c *gin.Context) {
	levelId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if levelId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	ent, err := dao.GetUserLevelById(int(levelId))

	//_, err := base.GetByCol("id", levelId, UserLevelEnt)
	// exist, err := base.GetByCol("levelId", levelId, UserLevelEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if ent == nil {
		StdErrResponse(c, ErrLevelNotExist)
		return
	}

	if enable == 1 {
		// 启用
		ent.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		ent.Status = base.StatusDeleted
	}

	// todo 操作日志

	_, err = base.UpdateByCol("id", levelId, ent, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteUserLevelList)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func init() {
}

type ShareSearchRequest struct {
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"pagesize"`
}

func (r *ShareSearchRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.Name != "" {
		condition["fileName like ?"] = "%" + r.Name + "%"
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

func ShareSearch(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ShareSearchRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Share, 0)
	var total int64

	if req.Name != "" {
		total, err = base.ListAndCountByCondition(&dao.Share{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
		if err != nil {
			log.Trace(requestId).Errorf("list failed:%v", err)
			StdErrResponse(c, ErrInternal)
			return
		}
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*ShareItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewShareItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "share/search.tpl", gin.H{
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

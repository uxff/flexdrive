package handler

import (
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
		JsonErr(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Share, 0)
	var total int64

	if req.Name != "" {
		total, err = base.ListAndCountByCondition(&dao.Share{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
		if err != nil {
			log.Trace(requestId).Errorf("list failed:%v", err)
			JsonErr(c, ErrInternal)
			return
		}
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*ShareItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewShareItemFromEnt(v))
	}

	JsonOk(c, gin.H{
		"total":     total,
		"page":      req.Page,
		"pagesize":  req.PageSize,
		"list":      resItems,
		"reqParam":  req,
		"paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

func ShareDetail(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	shareHash := c.Param("shareHash")
	if shareHash == "" {
		JsonErr(c, ErrInvalidParam)
		return
	}

	shareItem, err := dao.GetShareByShareHash(shareHash)
	if err != nil {
		log.Trace(requestId).Debugf("get shareHash(%s) error:%v", shareHash, err)
		JsonErr(c, ErrInternal)
		return
	}

	if shareItem == nil || shareItem.Status == base.StatusDeleted {
		JsonErrMsg(c, ErrItemNotExist, "分享的内容不存在或已删除")
		return
	}

	genShareOutPath(c, shareItem)

	JsonOk(c, shareItem)
}

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
	tplFuncMap["orderStatus"] = func(status int) string {
		return dao.OrderStatusMap[status]
	}
}

type OrderListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	UserEmail   string `form:"userEmail"`
	UserId      int    `form:"-"`
	Status      int    `form:"status"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *OrderListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.CreateStart != "" {
		condition["created>=?"] = r.CreateStart
	}

	if r.CreateEnd != "" {
		condition["created<=?"] = r.CreateEnd
	}

	if r.Status != 0 {
		condition["status=?"] = r.Status
	}

	if r.UserEmail != "" {
		userEnt, _ := dao.GetUserByEmail(r.UserEmail)
		if userEnt != nil && userEnt.Id > 0 {
			condition["userId=?"] = r.UserId
		} else {
			condition["userId=?"] = 0
		}
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type OrderItem struct {
	dao.Order
}

func NewOrderItemFromEnt(orderEnt *dao.Order) *OrderItem {
	return &OrderItem{
		Order: *orderEnt,
	}
}

func OrderList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &OrderListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Order, 0)
	total, err := base.ListAndCountByCondition(&dao.Order{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*OrderItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewOrderItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "order/list.tpl", gin.H{
		"LoginInfo":      getLoginInfo(c),
		"IsLogin":        isLoginIn(c),
		"total":          total,
		"page":           req.Page,
		"pagesize":       req.PageSize,
		"list":           resItems,
		"reqParam":       req,
		"orderStatusMap": dao.OrderStatusMap,
		"paginator":      paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

func OrderRefund(c *gin.Context) {
	orderId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if orderId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	orderEnt, err := dao.GetOrderById(int(orderId))
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if orderEnt == nil {
		StdErrResponse(c, ErrItemNotExist)
		return
	}

	if orderEnt.Status != dao.OrderStatusPaid {
		StdErrMsgResponse(c, ErrInternal, "该订单不允许退款")
		return
	}

	// todo call api for refund

	orderEnt.Status = dao.OrderStatusRefended

	//base.CacheDelByEntity("mgrLoginName", orderEnt.Email, orderEnt)

	_, err = base.UpdateByCol("id", orderId, orderEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteOrderList)
}

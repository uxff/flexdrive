package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func init() {
}

type OrderListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"name"`
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

	if r.Name != "" {
		condition["levelName like ?"] = "%" + r.Name + "%"
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type OrderItem struct {
	dao.Order
}

func NewOrderItemFromEnt(OrderEnt *dao.Order) *OrderItem {
	return &OrderItem{
		Order: *OrderEnt,
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
	var total int64

	if req.Name != "" {
		total, err = base.ListAndCountByCondition(&dao.Order{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
		if err != nil {
			log.Trace(requestId).Errorf("list failed:%v", err)
			StdErrResponse(c, ErrInternal)
			return
		}
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*OrderItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewOrderItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "order/list.tpl", gin.H{
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

func Mockpay(c *gin.Context) {
	orderIdStr := c.Query("orderId")
	if orderIdStr == "" {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	orderId, _ := strconv.Atoi(orderIdStr)

	orderInfo, err := dao.GetOrderById(orderId)
	if err != nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询错误")
		return
	}

	if orderInfo == nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询不到订单")
		return
	}

	loginInfo := getLoginInfo(c)
	if loginInfo.UserId != orderInfo.UserId {
		StdErrMsgResponse(c, ErrInvalidParam, "无权限操作该订单")
		return
	}

	verifyCode := fmt.Sprintf("%d", (time.Now().Unix()/79)%9999)
	outOrderNo := fmt.Sprintf("%d", (time.Now().UnixNano() / 777777))

	c.HTML(http.StatusOK, "order/list.tpl", gin.H{
		"LoginInfo":  getLoginInfo(c),
		"IsLogin":    isLoginIn(c),
		"Order":      orderInfo,
		"OutOrderNo": outOrderNo,
		"VerifyCode": verifyCode,
	})
}
func MockpayForm(c *gin.Context) {
	orderIdStr := c.Query("orderId")
	if orderIdStr == "" {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	orderId, _ := strconv.Atoi(orderIdStr)

	orderInfo, err := dao.GetOrderById(orderId)
	if err != nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询错误")
		return
	}

	if orderInfo == nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询不到订单")
		return
	}

	loginInfo := getLoginInfo(c)
	if loginInfo.UserId != orderInfo.UserId {
		StdErrMsgResponse(c, ErrInvalidParam, "无权限操作该订单")
		return
	}

	verifyCode := c.Query("verifyCode")
	if verifyCode == "" {
		StdErrMsgResponse(c, ErrInvalidParam, "验证码不对")
		return
	}

	verifyCodeExpected := fmt.Sprintf("%d", (time.Now().Unix()/79)%9999)
	if verifyCode != verifyCodeExpected {
		StdErrMsgResponse(c, ErrInvalidParam, "验证码不对，期望"+verifyCodeExpected)
		return
	}

	outOrderNo := c.Query("outOrderNo")
	if outOrderNo == "" {
		StdErrMsgResponse(c, ErrInvalidParam, "未收到第三方支付订单，支付失败")
		return
	}

	orderInfo.Remark = outOrderNo
	orderInfo.Status = dao.OrderStatusPaid
	orderInfo.UpdateById([]string{"remark", "status"})

	StdErrResponse(c, ErrSuccess)
}

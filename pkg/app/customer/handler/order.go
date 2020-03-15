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
	"github.com/uxff/flexdrive/pkg/utils"
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

	loginInfo := getLoginInfo(c)
	condition := req.ToCondition()
	condition["userId=?"] = loginInfo.UserId

	// 列表查询
	list := make([]*dao.Order, 0)
	var total int64

	total, err = base.ListAndCountByCondition(&dao.Order{}, condition, req.Page, req.PageSize, "", &list)
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

func OrderCreate(c *gin.Context) {

	c.HTML(http.StatusOK, "order/create.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
	})
}

func OrderCreateForm(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	levelIdStr := c.PostForm("level")
	if levelIdStr == "" {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	levelId, _ := strconv.Atoi(levelIdStr)

	levelInfo, err := dao.GetUserLevelById(levelId)
	if err != nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询等级错误")
		return
	}

	if levelInfo == nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询不到要购买的等级")
		return
	}

	loginInfo := getLoginInfo(c)

	orderInfo := &dao.Order{
		UserId:        loginInfo.UserId,
		OriginLevelId: loginInfo.UserEnt.LevelId,
		TotalAmount:   levelInfo.Price,
		PayAmount:     levelInfo.Price,
		AwardSpace:    levelInfo.QuotaSpace,
		AwardLevelId:  levelInfo.Id,
		LevelName:     levelInfo.Name,
		Status:        1,
	}

	_, err = base.Insert(orderInfo)
	if err != nil {
		log.Trace(requestId).Errorf("insert order error:%v", err)
		StdErrMsgResponse(c, ErrInvalidParam, "创建订单失败")
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/my/order/detail/"+strconv.Itoa(orderInfo.Id))
}

func OrderDetail(c *gin.Context) {

	orderIdStr := c.Param("orderId")
	if orderIdStr == "" {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	orderId, _ := strconv.Atoi(orderIdStr)

	orderInfo, err := dao.GetOrderById(orderId)
	if err != nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询订单错误")
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

	levelInfo, err := dao.GetUserLevelById(orderInfo.AwardLevelId)
	if err != nil {
		StdErrMsgResponse(c, ErrInvalidParam, "查询等级错误")
		return
	}

	if levelInfo == nil {
		StdErrMsgResponse(c, ErrInvalidParam, "等级不存在")
		return
	}

	c.HTML(http.StatusOK, "order/detail.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"Order":     orderInfo,
		"Level":     levelInfo,
	})
}

func Mockpay(c *gin.Context) {
	orderIdStr := c.Param("orderId")
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

	verifyCode := genVerifyCode(orderId)
	outOrderNo := fmt.Sprintf("%d", (time.Now().UnixNano() / 7777777))
	token := utils.Md5(outOrderNo + "/")

	c.HTML(http.StatusOK, "order/mockpay.tpl", gin.H{
		//"LoginInfo":  getLoginInfo(c),
		//"IsLogin":    isLoginIn(c),
		"Order":      orderInfo,
		"OutOrderNo": outOrderNo,
		"VerifyCode": verifyCode,
		"Token":      token,
		"User":       loginInfo.UserEnt,
	})
}

func MockpayForm(c *gin.Context) {
	orderIdStr := c.PostForm("orderId")
	if orderIdStr == "" {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	orderId, _ := strconv.Atoi(orderIdStr)

	verifyCode := c.PostForm("verifyCode")
	if !checkVerifyCode(orderId, verifyCode) {
		StdErrMsgResponse(c, ErrInvalidParam, "验证码不对")
		return
	}

	token := c.PostForm("token")
	if token == "" {
		StdErrMsgResponse(c, ErrInvalidParam, "表单token失效")
		return
	}

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

	outOrderNo := c.PostForm("outOrderNo")
	if outOrderNo == "" {
		StdErrMsgResponse(c, ErrInvalidParam, "未收到第三方支付订单，支付失败")
		return
	}

	orderInfo.Remark = outOrderNo
	orderInfo.Status = dao.OrderStatusPaid
	orderInfo.UpdateById([]string{"remark", "status"})

	// 用户升级
	loginInfo.UserEnt.QuotaSpace += orderInfo.AwardSpace
	loginInfo.UserEnt.TotalCharge += orderInfo.TotalAmount
	loginInfo.UserEnt.LevelId = orderInfo.AwardLevelId

	loginInfo.UserEnt.UpdateById([]string{"quotaSpace", "totalCharge", "levelId"})

	StdErrMsgResponse(c, ErrSuccess, "支付完成")
}

// 每秒都能产生不同的验证码 79秒内都可以正确验证
func genVerifyCode(orderId int) string {
	fixCode := int(time.Now().UnixNano()/1001) % 16
	longCode := utils.Md5(fmt.Sprintf("%d-%d", time.Now().Unix()/79, fixCode))

	return fmt.Sprintf("%x%s", fixCode, longCode[:3])
}

func checkVerifyCode(orderId int, v string) bool {
	if v == "" {
		return false
	}
	fixCode, _ := strconv.ParseInt(v[0:1], 16, 32)
	longCode := utils.Md5(fmt.Sprintf("%d-%d", time.Now().Unix()/79, fixCode))

	return v[1:] == longCode[:3]
}

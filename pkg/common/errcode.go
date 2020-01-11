package common

// 支付中心业务处理流程

const (
	// 支付中心内部错误
	Success              = "00000" // 操作成功
	InternalError        = "10000" // 系统内部错误 // 内部BUG级别错误,不应该出现的错误
	DatabaseReadError    = "10001" // 数据库读操作失败
	DatabaseWriteError   = "10002" // 数据库写操作失败
	UnknownStatusInDB    = "10003" // 数据库中的状态不合法
	GetMchConfigError    = "10004" // 获取商户配置失败
	SaveMchConfigError   = "10005" // 保存商户配置失败
	MchConfigInvalid     = "10006" // 商户配置格式错误
	MchNotExist          = "10007" // 商户不存在
	PayChannelNotSupport = "10008" // 配置的支付通道不被支持
	OrderNotExist        = "10009" // 订单不存在
	AppIdNotExist        = "10101" // appId不存在
	NoBizAccountConfig   = "10102" // bizAccount对应场景的bac不存在

	// 客户端错误
	SignError                   = "20001" // 签名验证失败
	SceneNotSupport             = "20002" // 不支持的支付场景
	ParamsError                 = "20003" // 客户端请求参数错误
	BizOrderNoNotExist          = "20004" // 商户订单号不存在
	BizOrderNoRepeated          = "20006" // 商户订单号重复
	MchIdNotExist               = "20007" // 商户号不存在
	OriginOrderNotAllowedRefund = "20008" // 原订单状态不允许退款
	OriginOrderNotAllowedClose  = "20009" // 不允许关闭订单
	BaseRequestParseError       = "20010" // 公共参数格式错误
	BaIdNotFound                = "20011" // baId格式非法或不存在
	SignTypeNotSupport          = "20012" // 不支持的签名类型
	BizDataParseError           = "20013" // 解析bizData失败
	OrderNoParamsEmpty          = "20014" // 没有传可用的订单号
	NotForwardOrder             = "20015" // 不是正向订单
	NotReverseOrder             = "20016" // 不是退款订单
	InvalidOrderNo              = "20017" // 订单号格式错误
	BizOrderNoNotSet            = "20018" // 退款时未传原业务方原支付订单号
	BizRefundOrderNoNotSet      = "20019" // 退款时未传业务方本次退款订单号

	// 通道方返回错误
	HttpRequestError      = "30001" // 向通道发起http请求失败
	ParseChannelDataError = "30002" // 解析通道返回的应答失败
	ChannelBizError       = "30003" // 通道操作失败
	ChannelRefundFailed   = "30004" // 通道发起退款失败
)

//var statusText = map[string]string{
//	// 支付中心内部错误
//	Success:              "操作成功",
//	InternalError:        "系统内部错误",
//	DatabaseReadError:    "数据库读操作失败",
//	DatabaseWriteError:   "数据库写操作失败",
//	UnknownStatusInDB:    "数据库中的状态不合法",
//	GetMchConfigError:    "获取商户配置失败",
//	SaveMchConfigError:   "保存商户配置失败",
//	MchConfigInvalid:     "商户配置格式错误",
//	MchNotExist:          "商户不存在",
//	PayChannelNotSupport: "配置的支付通道不被支持",
//	OrderNotExist:        "订单不存在",
//	AppIdNotExist:        "appId不存在",
//	NoBizAccountConfig:   "bizAccount对应场景的bac不存在",
//
//	// 客户端错误
//	SignError:                   "签名验证失败",
//	SceneNotSupport:             "不支持的支付场景",
//	ParamsError:                 "客户端请求参数错误",
//	BizOrderNoNotExist:          "商户订单号不存在",
//	BizOrderNoRepeated:          "商户订单号重复",
//	MchIdNotExist:               "商户号不存在",
//	OriginOrderNotAllowedRefund: "原订单状态不允许退款",
//	OriginOrderNotAllowedClose:  "不允许关闭订单",
//	BaseRequestParseError:       "公共参数格式错误",
//	BaIdNotFound:                "baId格式非法或不存在",
//	SignTypeNotSupport:          "不支持的签名类型",
//	BizDataParseError:           "解析bizData失败",
//	OrderNoParamsEmpty:          "没有传可用的订单号",
//	NotForwardOrder:             "不是正向订单",
//	NotReverseOrder:             "不是退款订单",
//	InvalidOrderNo:              "订单号格式错误",
//	BizOrderNoNotSet:            "退款时未传原业务方原支付订单号",
//	BizRefundOrderNoNotSet:      "退款时未传业务方本次退款订单号",
//
//	// 通道方返回错误
//	HttpRequestError:      "向通道发起http请求失败",
//	ParseChannelDataError: "解析通道返回的应答失败",
//	ChannelBizError:       "通道操作失败",
//	ChannelRefundFailed:   "通道发起退款失败",
//}

const unknownStatus = "Unknown Status"

//func StatusText(status string) string {
//	v, ok := statusText[status]
//	if ok {
//		return v
//	}
//	return unknownStatus
//}

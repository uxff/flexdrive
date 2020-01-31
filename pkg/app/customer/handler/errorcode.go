package handler

const (
	// 成功
	ErrSuccess = "0"

	// 错误
	ErrNotLogin       = "10"
	ErrLoginExpired   = "11"
	ErrUserNotExist   = "12"
	ErrInvalidPass    = "13"
	ErrNoPermit       = "14"
	ErrInvalidParam   = "15"
	ErrInvalidCaptcha = "16"
	ErrNameDuplicate  = "17"
	ErrUserDisabled   = "18"
	ErrInternal       = "101"
	ErrLevelNotExist  = "21"
	ErrLevelDisabled  = "22"
	ErrEmailDuplicate = "23"
	ErrItemNotExist   = "24"
)

var errCodeMap = map[string]string{
	ErrSuccess:        "操作成功",
	ErrNotLogin:       "尚未登录，需要登录后才能操作",
	ErrLoginExpired:   "登录状态已过期，需要重新登录",
	ErrUserNotExist:   "账号不存在",
	ErrUserDisabled:   "账号被禁用",
	ErrInvalidPass:    "账号错误或密码错误",
	ErrNoPermit:       "没有权限",
	ErrInvalidParam:   "参数错误",
	ErrInvalidCaptcha: "验证码错误",
	ErrNameDuplicate:  "名称重复，请更换名称后重试",
	ErrInternal:       "系统内部错误，请稍后再试",
	ErrLevelDisabled:  "会员等级被禁用",
	ErrLevelNotExist:  "会员等级不存在",
	ErrEmailDuplicate: "邮箱重复",
	ErrItemNotExist:   "该条目不存在",
}

func CodeToMessage(code string) string {
	return errCodeMap[code]
}

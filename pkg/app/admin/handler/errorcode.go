package handler

const (
	// 成功
	ErrSuccess = "0"

	// 错误
	ErrNotLogin       = "10"
	ErrLoginExpired   = "11"
	ErrMgrNotExist    = "12"
	ErrInvalidPass    = "13"
	ErrNoPermit       = "14"
	ErrInvalidParam   = "15"
	ErrInvalidCaptcha = "16"
	ErrNameDuplicate  = "17"
	ErrMgrDisabled    = "18"
	ErrInternal       = "101"
)

var errCodeMap = map[string]string{
	ErrSuccess: "成功",
}

func codeToMessage(code string) string {
	return errCodeMap[code]
}

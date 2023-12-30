package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils"
)

// cookie中使用
const (
	CookieKeyAuth = "ua"
	CookieKeySalt = "TmhMbU52YlM1amJp"
)

// 代码内使用 http协议中不可见
const (
	CtxKeyCua          = "_cua"
	CtxKeyRequestId    = "_requestId"
	CtxKeyURI          = "_uri"
	CookieKeyCaptchaId = "_captchaId"
)

const (
	cuaTimeDiv = 79
)

const (
	LoginCookieExpire = 3600 * 24 * 365 // 365天
)

// 登录关键信息 Customer User Auth Token
type CuaToken struct {
	UserId  int
	LoginAt int
	Sign    string

	UserEnt *dao.User
}

func (t *CuaToken) ToString() string {
	return fmt.Sprintf("%d.%d.%s", t.LoginAt/cuaTimeDiv, t.UserId, t.Sign)
}

func (t *CuaToken) FromString(str string) {
	t.UserId = 0
	cols := strings.Split(str, ".")
	if len(cols) != 3 {
		return
	}
	t.UserId, _ = strconv.Atoi(cols[1])
	t.LoginAt, _ = strconv.Atoi(cols[0])
	t.LoginAt *= cuaTimeDiv
	t.Sign = cols[2]
}

func (t *CuaToken) MakeSign() string {

	enc := md5.New()
	enc.Write([]byte(fmt.Sprintf("%d.%d", t.UserId, t.LoginAt/cuaTimeDiv) + CookieKeySalt))

	return hex.EncodeToString(enc.Sum(nil))
}

// func decodeCuaFromToken(cuaTokenStr string) (g *CuaToken, err error) {

// 	g = &CuaToken{}
// 	g.FromString(cuaTokenStr)

// 	if g.UserId <= 0 {
// 		return nil, errors.New("cuaToken has no userId")
// 	}

// 	if g.MakeSign() != g.Sign {
// 		return nil, errors.New("cuaToken sign not expected")
// 	}

// 	return g, nil
// }

// func genCuaFromUserEnt(userEnt *dao.User) (g *CuaToken, cuaTokenStr, sign string, err error) {
// 	g = &CuaToken{
// 		UserId:  userEnt.Id,
// 		LoginAt: int(userEnt.LastLoginAt.Unix()),
// 		UserEnt: userEnt,
// 	}

// 	// toString 之前必须把签名赋值给g.Sign
// 	g.Sign = g.MakeSign()

// 	return g, g.ToString(), g.Sign, nil
// }

func decodeGpaFromJwtClaim(claim jwt.MapClaims) (g *CuaToken, err error) {

	g = &CuaToken{}
	g.FromString(claim["gpa"].(string))

	if g.UserId <= 0 {
		return nil, errors.New("CuaToken has no mid")
	}

	return g, nil
}

func genJwtClaimFromUserEnt(userEnt *dao.User) map[string]interface{} {
	g := &CuaToken{
		UserId:  userEnt.Id,
		LoginAt: int(userEnt.LastLoginAt.Unix()),
		UserEnt: userEnt,
	}
	return map[string]interface{}{
		"gpa": g.ToString(),
	}
}

func TraceMiddleWare(c *gin.Context) {
	// 埋入requestId
	requestId := utils.NewRandomHex(16)
	uri := c.Request.URL.String()
	c.Set(CtxKeyRequestId, requestId)
	c.Set(CtxKeyURI, uri)

	//rawBody, _ := httputil.DumpRequest(c.Request, true)
	//log.Trace(requestId).Debugf("原始请求体：%s", rawBody)

	// detect user from cookie // necessary?
	cuaToken, err := verifyFromCookie(c)
	if err != nil {
		log.Trace(c.GetString(CtxKeyRequestId)).Warnf("illegal cuaToken , reject request, error:%v cuatoken:%+v", err, cuaToken)
	} else {
		if cuaToken != nil {
			userEnt, err := dao.GetUserById(cuaToken.UserId)
			if err != nil {
				log.Warnf("query by userid:%d failed:%v", cuaToken.UserId, err)
			}
			if userEnt != nil {
				// 成功将用户实体注入到登录信息回话
				cuaToken.UserEnt = userEnt
				c.Set(CtxKeyCua, cuaToken)
			}
		}
	}

	c.Next()
}

// 对每个请求添加全局requestId，放到gin.Context里。后面的日志里尽量加上，方便追踪问题
// 所有交易相关接口调用前的认证中间件
func AuthMiddleWare(c *gin.Context) {
	// 验证cookie签名是否合法
	cuaToken := getLoginInfo(c)
	if cuaToken == nil {
		log.Warnf("尚未登录")
		JsonErr(c, ErrNotLogin)
		c.Abort()
		return
	}

	if cuaToken.LoginAt < int(time.Now().Add(-time.Hour*24*7).Unix()) {
		ClearLogin(c)
		JsonErr(c, ErrLoginExpired)
		c.Abort()
		return
	}

	// 判断账号是否已被禁用
	if cuaToken.UserEnt.Status != base.StatusNormal {
		log.Warnf("登陆账号(%d)已被禁用", cuaToken.UserId)
		JsonErr(c, ErrUserDisabled)
		c.Abort()
		return
	}

	c.Next()
}

// // 接口调用未出错时，标准输出必须调用的接口
// func StdResponse(c *gin.Context, code string, biz interface{}) {
// 	StdResponseJson(c, code, "", biz)
// }

// // 接口调用出错时，标准输出必须调用的接口
// func StdErrResponse(c *gin.Context, code string) {
// 	errMsg := CodeToMessage(code)
// 	c.HTML(http.StatusOK, "common/error.tpl", gin.H{
// 		"LoginInfo": getLoginInfo(c),
// 		"IsLogin":   isLoginIn(c),
// 		"errMsg":    errMsg,
// 	})
// 	//StdResponseJson(c, code, "", "")
// }
// func StdErrMsgResponse(c *gin.Context, code string, errMsg string) {
// 	if errMsg == "" {
// 		errMsg = CodeToMessage(code)
// 	}
// 	c.HTML(http.StatusOK, "common/error.tpl", gin.H{
// 		"LoginInfo": getLoginInfo(c),
// 		"IsLogin":   isLoginIn(c),
// 		"errMsg":    errMsg,
// 	})
// 	//StdResponseJson(c, code, "", "")
// }

// func StdResponseJson(c *gin.Context, code, msg string, data interface{}) {
// 	requestId := c.GetString(CtxKeyRequestId)

// 	codeMsg := CodeToMessage(code)
// 	if msg != "" {
// 		codeMsg += "(" + msg + ")"
// 	}

// 	resp := gin.H{
// 		"errcode":   code,
// 		"errmsg":    codeMsg,
// 		"errlevel":  "alert",
// 		"result":    data,
// 		"requestId": requestId,
// 	}

// 	c.JSON(http.StatusOK, resp)

// 	b, _ := json.Marshal(resp)
// 	log.Trace(requestId).Debugf("==========DEBUG - URI:%s 应答：%+s", c.GetString(CtxKeyURI), b)

// 	// log.Trace(requestId).Infof("URI:%s 应答：%+v", c.GetString(CtxKeyURI), resp)
// }

// 接口调用未出错时，标准输出必须调用的接口
func JsonOk(c *gin.Context, data interface{}) {
	JsonResponse(c, ErrSuccess, "", data)
}

// 接口调用出错时，标准输出必须调用的接口
func JsonErr(c *gin.Context, code string) {
	JsonResponse(c, code, CodeToMessage(code), "")
}
func JsonErrMsg(c *gin.Context, code string, errMsg string) {
	JsonResponse(c, code, errMsg, "")
}

func JsonResponse(c *gin.Context, code, msg string, data interface{}) {
	requestId := c.GetString(CtxKeyRequestId)

	codeMsg := CodeToMessage(code)
	if msg != "" {
		codeMsg += "(" + msg + ")"
	}

	resp := gin.H{
		"errcode":   code,
		"errmsg":    codeMsg,
		"errlevel":  "alert",
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"result":    data,
		"requestId": requestId,
	}

	c.JSON(http.StatusOK, resp)

	b, _ := json.Marshal(resp)
	log.Trace(requestId).Debugf("==========DEBUG - URI:%s 应答：%+s", c.GetString(CtxKeyURI), b)

	// log.Trace(requestId).Infof("URI:%s 应答：%+v", c.GetString(CtxKeyURI), resp)
}

// 此方法必须提前验证cookie 就是前文必须调用过verifyFromCookie，此方法才有效
func getLoginInfo(c *gin.Context) *CuaToken {
	loginInfoIf, _ := c.Get(CtxKeyCua)
	loginInfo, ok := loginInfoIf.(*CuaToken)
	if !ok {
		log.Warnf("cuaToken not exist, invalid type")
	}
	if loginInfo == nil {
		log.Warnf("cuaToken not exist, empty value")
	}
	return loginInfo
}

// 验证cookie合法性 并返回有效的登录信息
func verifyFromCookie(c *gin.Context) (*CuaToken, error) {
	// gopay admin token
	CuaTokenStr, err := c.Cookie(CookieKeyAuth)
	if CuaTokenStr == "" {
		return nil, err
	}

	// // gopay admin sign
	// gpaSignStr, err := c.Cookie(CookieKeySign)
	// if err != nil {
	// 	return nil, err
	// }

	jwtToken, err := jwt.Parse(CuaTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(CookieKeySalt), nil
	})

	if err != nil || jwtToken == nil {
		log.Warnf("get jwtToken failed, error:%v", err)
		return nil, err
	}

	// log.Debugf("jwtToken:%+v", jwtToken)

	CuaToken, err := decodeGpaFromJwtClaim(jwtToken.Claims.(jwt.MapClaims))
	if err != nil {
		log.Warnf("get CuaToken from jwt failed, error:%v", err)
		return nil, err
	}

	return CuaToken, nil
}
func isLoginIn(c *gin.Context) bool {
	return getLoginInfo(c) != nil
}

// // 验证cookie合法性 并返回有效的登录信息
// func verifyFromCookie(c *gin.Context) (*CuaToken, error) {
// 	// customer auth token
// 	cuaTokenStr, err := c.Cookie(CookieKeyAuth)
// 	if cuaTokenStr == "" {
// 		return nil, err
// 	}

// 	cuaToken, err := decodeCuaFromToken(cuaTokenStr)
// 	if err != nil {
// 		log.Warnf("cookie has invalid sign, error:%v", err)
// 		return nil, err
// 	}

// 	return cuaToken, nil
// }

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
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils"
)

// cookie中使用
const (
	CookieKeyGpa  = "t"
	CookieKeySign = "s"
	CookieKeySalt = "TmhMbU52YlM1amJp"
)

// 代码内使用 http协议中不可见
const (
	CtxKeyGpa          = "_gpa"
	CtxKeyRequestId    = "_requestId"
	CtxKeyURI          = "_uri"
	CookieKeyCaptchaId = "_captchaId"
)

const (
	LoginCookieExpire = 3600 * 24 * 365 // 365天
)

// 后台登录关键信息
type GpaToken struct {
	Mid int
	//Name     string // 没有值
	RoleId int
	//RoleName string // 没有值
	LoginAt int

	MgrEnt *dao.Manager
}

func (t *GpaToken) ToString() string {
	return fmt.Sprintf("%d.%d.%d", t.Mid, t.RoleId, t.LoginAt)
}

func (t *GpaToken) FromString(str string) {
	t.Mid = 0
	cols := strings.Split(str, ".")
	if len(cols) != 3 {
		return
	}
	t.Mid, _ = strconv.Atoi(cols[0])
	t.RoleId, _ = strconv.Atoi(cols[1])
	t.LoginAt, _ = strconv.Atoi(cols[2])
}

func (t *GpaToken) MakeSign() string {

	enc := md5.New()
	enc.Write([]byte(t.ToString() + CookieKeySalt))

	return hex.EncodeToString(enc.Sum(nil))
}

func decodeGpaFromToken(gpaTokenStr string, sign string) (g *GpaToken, err error) {

	g = &GpaToken{}
	g.FromString(gpaTokenStr)

	if g.Mid <= 0 {
		return nil, errors.New("gpatoken has no mid")
	}

	if g.MakeSign() != sign {
		return nil, errors.New("gpatoken sign not expected")
	}

	return g, nil
}

func genGpaFromMgrEnt(mgrEnt *dao.Manager) (g *GpaToken, gpaTokenStr, sign string, err error) {
	g = &GpaToken{
		Mid: mgrEnt.Id,
		//Name:   mgrEnt.Name,
		RoleId: mgrEnt.RoleId,
		//RoleName: 	 mgrEnt.RoleName,
		LoginAt: int(mgrEnt.LastLoginAt.Unix()),
	}

	return g, g.ToString(), g.MakeSign(), nil
}

func TraceMiddleWare(c *gin.Context) {
	// 埋入requestId
	requestId := utils.NewRandomHex(16)
	uri := c.Request.URL.String()
	c.Set(CtxKeyRequestId, requestId)
	c.Set(CtxKeyURI, uri)

	//rawBody, _ := httputil.DumpRequest(c.Request, true)
	//log.Trace(requestId).Debugf("原始请求体：%s", rawBody)

	c.Next()
}

// 对每个请求添加全局requestId，放到gin.Context里。后面的日志里尽量加上，方便追踪问题
// 所有交易相关接口调用前的认证中间件
func AuthMiddleWare(c *gin.Context) {
	// 验证cookie签名是否合法
	gpaToken, err := verifyFromCookie(c)
	if err != nil {
		log.Trace(c.GetString(CtxKeyRequestId)).Warnf("illegal gpaToken , reject request, error:%v gpatoken:%+v", err, gpaToken)
		// c.SetCookie(CookieKeyGpa, "", -1, "", "", true, false)
		StdErrResponse(c, ErrNotLogin)
		c.Abort()
		return
	}

	if gpaToken.LoginAt+LoginCookieExpire < int(time.Now().Unix()) {
		// ClearLogin(c)
		StdErrResponse(c, ErrLoginExpired)
		c.Abort()
		return
	}

	//mgrEnt := &dao.Manager{}
	mgrEnt, err := dao.GetManagerById(gpaToken.Mid)
	if err != nil {
		log.Errorf("query by mid:%d failed:%v", gpaToken.Mid, err)
		StdErrResponse(c, ErrMgrNotExist)
		c.Abort()
		return
	}
	if mgrEnt == nil {
		log.Warnf("登陆账号%d不存在", gpaToken.Mid)
		StdErrResponse(c, ErrMgrNotExist)
		c.Abort()
		return
	}

	gpaToken.RoleId = mgrEnt.RoleId
	gpaToken.MgrEnt = mgrEnt
	//gpaToken.RoleName = mgrEnt.RoleName
	// gpaToken.IsSuper = mgrEnt.IsSuper()
	// gpaToken.IsSuperRole = mgrEnt.IsSuperRole()

	c.Set(CtxKeyGpa, gpaToken)

	// 判断账号是否已被禁用
	if mgrEnt.Status != base.StatusNormal {
		log.Warnf("登陆账号(%d)已被禁用", gpaToken.Mid)
		StdErrResponse(c, ErrMgrDisabled)
		c.Abort()
		return
	}

	c.Next()
}

func RbacAuthMiddleWare(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	gpaToken := getLoginInfo(c)
	if gpaToken.MgrEnt.IsSuper == 1 {
		c.Next()
		return
	}

	// 基于角色鉴权
	//roleEnt, err := roles.QueryByRid(gpaToken.RoleId)
	roleEnt, err := dao.GetRoleById(gpaToken.MgrEnt.RoleId)
	if err != nil || roleEnt == nil || roleEnt.Status != base.StatusNormal {
		log.Trace(requestId).Errorf("get role(%d) failed:%v", gpaToken.RoleId, err)
		StdErrResponse(c, ErrRoleNotExist)
		c.Abort()
		return
	}
	//
	if roleEnt.IsSuper() {
		// 超级管理忽略权限
		c.Next()
		return
	}
	//
	if !roleEnt.Permit.CheckRouteAccessable(c.Request.RequestURI) {
		log.Trace(requestId).Errorf("no access, roleid:%d route:%s", gpaToken.RoleId, c.Request.RequestURI)
		StdErrResponse(c, ErrNoPermit)
		c.Abort()
		return
	}
	log.Trace(requestId).Debugf("access allowed")

	c.Next()
}

// 接口调用未出错时，标准输出必须调用的接口
func StdResponse(c *gin.Context, code string, biz interface{}) {
	StdResponseJson(c, code, "", biz)
}

// 接口调用出错时，标准输出必须调用的接口
func StdErrResponse(c *gin.Context, code string) {
	errMsg := CodeToMessage(code)
	c.HTML(http.StatusOK, "common/error.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"errMsg":    errMsg,
	})
	//StdResponseJson(c, code, "", "")
}
func StdErrMsgResponse(c *gin.Context, code string, errMsg string) {
	if errMsg == "" {
		errMsg = CodeToMessage(code)
	}
	c.HTML(http.StatusOK, "common/error.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"errMsg":    errMsg,
	})
	//StdResponseJson(c, code, "", "")
}

func StdResponseJson(c *gin.Context, code, msg string, data interface{}) {
	requestId := c.GetString(CtxKeyRequestId)

	codeMsg := CodeToMessage(code)
	if msg != "" {
		codeMsg += "(" + msg + ")"
	}

	resp := gin.H{
		"errcode":   code,
		"errmsg":    codeMsg,
		"errlevel":  "alert",
		"result":    data,
		"requestId": requestId,
	}

	c.JSON(http.StatusOK, resp)

	b, _ := json.Marshal(resp)
	log.Trace(requestId).Debugf("==========DEBUG - URI:%s 应答：%+s", c.GetString(CtxKeyURI), b)

	log.Trace(requestId).Infof("URI:%s 应答：%+v", c.GetString(CtxKeyURI), resp)
}

// 此方法必须提前验证cookie 就是前文必须调用过verifyFromCookie，此方法才有效
func getLoginInfo(c *gin.Context) *GpaToken {
	loginInfoIf, _ := c.Get(CtxKeyGpa)
	loginInfo, ok := loginInfoIf.(*GpaToken)
	if !ok {
		log.Warnf("gpatoken not exist, invalid type")
	}
	if loginInfo == nil {
		log.Warnf("gpatoken not exist, empty value")
	}
	return loginInfo
}

func isLoginIn(c *gin.Context) bool {
	return getLoginInfo(c) != nil
}

// 验证cookie合法性 并返回有效的登录信息
func verifyFromCookie(c *gin.Context) (*GpaToken, error) {
	// gopay admin token
	gpaTokenStr, err := c.Cookie(CookieKeyGpa)
	if gpaTokenStr == "" {
		return nil, err
	}

	// gopay admin sign
	gpaSignStr, err := c.Cookie(CookieKeySign)
	if err != nil {
		return nil, err
	}
	gpaToken, err := decodeGpaFromToken(gpaTokenStr, gpaSignStr)
	if err != nil {
		log.Warnf("cookie has invalid sign, error:%v", err)
		return nil, err
	}

	return gpaToken, nil
}

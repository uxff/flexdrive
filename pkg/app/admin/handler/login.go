package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"net/http"
	"time"
)

type LoginRequst struct {
	LoginName string `form:"loginName"`
	Pwd       string `form:"pwd"`
	Captcha   string `form:"captcha"`
}

type LoginResponse struct {
	LoginName string `json:"loginName"`
	Mid       int    `json:"mid"`
}

func Login(c *gin.Context) {
	//c.HTML(http.StatusOK, "pkg/app/admin/view/login/login.tpl", gin.H{})
	c.HTML(http.StatusOK, "login/login.tpl", gin.H{
		"path": "login",
	})
}

// 提交登录的处理
func LoginForm(c *gin.Context) {
	// 如果已经登录 则跳到成功页
	//gpaToken, _ := verifyFromCookie(c)
	//if gpaToken == nil {
	//	// 未登录
	//}

	// 参数是否正确
	req := &LoginRequst{}
	err := c.ShouldBind(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 验证码是否正确
	if !VerifyCaptcha(c, req.Captcha) {
		StdErrResponse(c, ErrInvalidCaptcha)
		return
	}

	mgrEnt := &dao.Manager{}

	err = mgrEnt.GetByName(req.LoginName)
	if err != nil {
		log.Errorf("query by mgrLoginName:%s failed:%v", req.LoginName, err)
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if mgrEnt.Id == 0 {
		log.Warnf("mgrLoginName:%s not exist, verify failed", req.LoginName)
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if mgrEnt.IsPwdValid(req.Pwd) {
		log.Warnf("mgr pwd not matched, verify failed. mgrLoginName:%s", req.LoginName)
		StdErrResponse(c, ErrInvalidPass)
	}

	//// 密码是否正确
	//mgrEnt := managermodel.VerifyPwd(req.LoginName, req.Pwd)
	//if mgrEnt == nil {
	//	StdErrResponse(c, errcodes.LoginFailed)
	//	return
	//}

	// 判断账号是否已被禁用
	if mgrEnt.Status != base.StatusNormal {
		StdErrResponse(c, ErrMgrDisabled)
		return
	}

	// 登录成功 种下cookie
	AcceptLogin(c, mgrEnt)

	StdResponse(c, ErrSuccess, LoginResponse{
		LoginName: mgrEnt.Name,
		Mid:       mgrEnt.Id,
	})
}

//// 获取一些全局配置 比如登录信息 菜单列表
//func GetAppConfig(c *gin.Context) {
//	gpaToken, _ := verifyFromCookie(c)
//
//	roleId := gpaToken.RoleId
//
//	roleEnt := &roles.Role{}
//	_, err := base.GetByColWithCache("rid", roleId, roleEnt)
//	if err != nil {
//		log.Errorf("get role(%d) menu failed:%v", roleId, err)
//		StdErrResponse(c, ErrInternal)
//		return
//	}
//
//	if !roleEnt.IsSuper() {
//		roleMenu := rbac.GetMenuByRoleEnt(roleEnt)
//		StdResponse(c, ErrSuccess, gin.H{
//			"loginInfo": gpaToken,
//			"nav":       menu.GetAllMenu(),
//			"rolenav":   roleMenu,
//			"access":    rbac.GetAccessMenuByRoleEnt(roleEnt),
//		})
//		return
//	}
//
//	StdResponse(c, ErrSuccess, gin.H{
//		"loginInfo": gpaToken,
//		"nav":       menu.GetAllMenu(),
//	})
//}

func Logout(c *gin.Context) {
	ClearLogin(c)
	StdResponse(c, ErrSuccess, nil)
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

// 受理登录
func AcceptLogin(c *gin.Context, mgrEnt *dao.Manager) {
	mgrEnt.LastLoginIp = c.Request.Header.Get("X-Real-IP")
	mgrEnt.LastLoginAt = time.Now() //util.JsonTime(time.Now()) // time.Now().Format("2006-01-02 15:04:05")

	_, tokenStr, sign, err := genGpaFromMgrEnt(mgrEnt)
	if err != nil {
		log.Errorf("gen gpatoken failed:%v", err)
		return
	}
	c.SetCookie(CookieKeyGpa, tokenStr, 3600*24*7, "", "", true, false)
	c.SetCookie(CookieKeySign, sign, 3600*24*7, "", "", true, false)

	// record login
	//go managermodel.RecordLoginStatus(mgrEnt)
}

func ClearLogin(c *gin.Context) {
	c.SetCookie(CookieKeyGpa, "", -1, "", "", true, false)
	c.SetCookie(CookieKeySign, "", -1, "", "", true, false)
}

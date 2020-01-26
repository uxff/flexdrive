package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

type LoginRequst struct {
	Email   string `form:"email"`
	Pwd     string `form:"password"`
	Captcha string `form:"captcha"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Mid   int    `json:"mid"`
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login/login.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
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

	mgrEnt, err := dao.GetManagerByEmail(req.Email)
	if err != nil {
		log.Errorf("query by email:%s failed:%v", req.Email, err)
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if mgrEnt.Id == 0 {
		log.Warnf("email:%s not exist, verify failed", req.Email)
		StdErrResponse(c, ErrMgrNotExist)
		return
	}

	if !mgrEnt.IsPwdValid(req.Pwd) {
		log.Warnf("mgr pwd not matched, verify failed. email:%s", req.Email)
		StdErrResponse(c, ErrInvalidPass)
		return
	}

	// 判断账号是否已被禁用
	if mgrEnt.Status != base.StatusNormal {
		StdErrResponse(c, ErrMgrDisabled)
		return
	}

	// 登录成功 种下cookie
	AcceptLogin(c, mgrEnt)

	c.Redirect(http.StatusMovedPermanently, RouteHome)
	//StdResponse(c, ErrSuccess, "/")
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
	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteLogin)
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
	c.SetCookie(CookieKeyGpa, tokenStr, 3600*24*7, "", "", false, false)
	c.SetCookie(CookieKeySign, sign, 3600*24*7, "", "", false, false)

	// record login
	//go managermodel.RecordLoginStatus(mgrEnt)
}

func ClearLogin(c *gin.Context) {
	c.SetCookie(CookieKeyGpa, "", -1, "", "", false, false)
	c.SetCookie(CookieKeySign, "", -1, "", "", false, false)
}

package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

type LoginRequest struct {
	Email   string `form:"email"`
	Pwd     string `form:"password"`
	Captcha string `form:"captcha"`
}

type LoginResponse struct {
	Email  string `json:"email"`
	UserId int    `json:"mid"`
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
	//cuaToken, _ := verifyFromCookie(c)
	//if cuaToken == nil {
	//	// 未登录
	//}

	// 参数是否正确
	req := &LoginRequest{}
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

	userEnt, err := dao.GetUserByEmail(req.Email)
	if err != nil {
		log.Errorf("query by email:%s failed:%v", req.Email, err)
		StdErrResponse(c, ErrUserNotExist)
		return
	}

	if userEnt.Id == 0 {
		log.Warnf("email:%s not exist, verify failed", req.Email)
		StdErrResponse(c, ErrUserNotExist)
		return
	}

	if !userEnt.IsPwdValid(req.Pwd) {
		log.Warnf("mgr pwd not matched, verify failed. email:%s", req.Email)
		StdErrResponse(c, ErrInvalidPass)
		return
	}

	// 判断账号是否已被禁用
	if userEnt.Status != base.StatusNormal {
		StdErrResponse(c, ErrUserDisabled)
		return
	}

	// 登录成功 种下cookie
	AcceptLogin(c, userEnt)

	c.Redirect(http.StatusMovedPermanently, RouteHome)
	//StdResponse(c, ErrSuccess, "/")
}

//// 获取一些全局配置 比如登录信息 菜单列表
//func GetAppConfig(c *gin.Context) {
//	cuaToken, _ := verifyFromCookie(c)
//
//	roleId := cuaToken.RoleId
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
//			"loginInfo": cuaToken,
//			"nav":       menu.GetAllMenu(),
//			"rolenav":   roleMenu,
//			"access":    rbac.GetAccessMenuByRoleEnt(roleEnt),
//		})
//		return
//	}
//
//	StdResponse(c, ErrSuccess, gin.H{
//		"loginInfo": cuaToken,
//		"nav":       menu.GetAllMenu(),
//	})
//}

func Logout(c *gin.Context) {
	ClearLogin(c)
	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteLogin)
}

// 受理登录
func AcceptLogin(c *gin.Context, userEnt *dao.User) {
	userEnt.LastLoginIp = c.Request.Header.Get("X-Real-IP")
	userEnt.LastLoginAt = time.Now() //util.JsonTime(time.Now()) // time.Now().Format("2006-01-02 15:04:05")

	_, tokenStr, _, err := genCuaFromUserEnt(userEnt)
	if err != nil {
		log.Errorf("gen gpatoken failed:%v", err)
		return
	}
	c.SetCookie(CookieKeyAuth, tokenStr, 3600*24*7, "", "", http.SameSiteDefaultMode, false, false)

	// record login
	//go managermodel.RecordLoginStatus(userEnt)
}

func ClearLogin(c *gin.Context) {
	c.SetCookie(CookieKeyAuth, "", -1, "", "", http.SameSiteDefaultMode, false, false)
}

func ChangePwd(c *gin.Context) {
	c.HTML(http.StatusOK, "login/changepwd.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
	})
}

type ChangePwdRequest struct {
	Email  string `form:"email"`
	Oldpwd string `form:"oldpwd"`
	Newpwd string `form:"newpwd"`
}

// 修改自己的密码
func ChangePwdForm(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	req := &ChangePwdRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	loginEnt := getLoginInfo(c)
	if loginEnt == nil || loginEnt.UserId <= 0 {
		StdErrResponse(c, ErrNotLogin)
		return
	}

	userEnt, err := dao.GetUserById(loginEnt.UserId)
	if err != nil {
		StdErrResponse(c, ErrUserNotExist)
		return
	}

	if req.Email != userEnt.Email {
		StdErrMsgResponse(c, ErrInvalidParam, "请提交自己的邮箱")
		return
	}

	if !userEnt.IsPwdValid(req.Oldpwd) {
		StdErrResponse(c, ErrInvalidPass)
		return
	}

	userEnt.SetPwd(req.Newpwd)

	// _, err = base.UpdateByCol("mid", loginEnt.UserId, mgrDbEnt, []string{"mgrPwd"})
	err = userEnt.UpdateById([]string{"pwd"})

	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteHome)
}

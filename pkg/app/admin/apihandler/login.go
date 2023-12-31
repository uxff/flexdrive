package apihandler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

// func Login(c *gin.Context) {
// 	c.HTML(http.StatusOK, "login/login.tpl", gin.H{
// 		"LoginInfo": getLoginInfo(c),
// 		"IsLogin":   isLoginIn(c),
// 	})
// }

type LoginRequest struct {
	Email   string `form:"email"`
	Pwd     string `form:"password"`
	Captcha string `form:"captcha"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Mid   int    `json:"mid"`
}

// 提交登录的处理
func LoginForm(c *gin.Context) {

	// 参数是否正确
	req := &LoginRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		JsonErr(c, ErrInvalidParam)
		return
	}

	// 验证码是否正确
	if !VerifyCaptcha(c, req.Captcha) {
		JsonErr(c, ErrInvalidCaptcha)
		return
	}

	// 查找账号
	mgrEnt, err := dao.GetManagerByEmail(req.Email)
	if err != nil {
		log.Errorf("query by email:%s failed:%v", req.Email, err)
		JsonErr(c, ErrMgrNotExist)
		return
	}

	if mgrEnt.Id == 0 {
		log.Warnf("email:%s not exist, verify failed", req.Email)
		JsonErr(c, ErrMgrNotExist)
		return
	}

	// 验证密码
	if !mgrEnt.IsPwdValid(req.Pwd) {
		log.Warnf("mgr pwd not matched, verify failed. email:%s", req.Email)
		JsonErr(c, ErrInvalidPass)
		return
	}

	// 判断账号是否已被禁用
	if mgrEnt.Status != base.StatusNormal {
		JsonErr(c, ErrMgrDisabled)
		return
	}

	// 登录成功 种下cookie
	AcceptLogin(c, mgrEnt)

	JsonOk(c, gin.H{
		"date": time.Now(),
	})
	// c.Redirect(http.StatusMovedPermanently, RouteHome)
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
	JsonOk(c, gin.H{
		"date": time.Now(),
	})
	// c.Redirect(http.StatusMovedPermanently, RouteLogin)
}

// 受理登录
func AcceptLogin(c *gin.Context, mgrEnt *dao.Manager) error {
	mgrEnt.LastLoginIp = c.Request.Header.Get("X-Real-IP")
	mgrEnt.LastLoginAt = time.Now() //util.JsonTime(time.Now()) // time.Now().Format("2006-01-02 15:04:05")

	// 生成签名
	// gpaToken, _, _, err := genGpaFromMgrEnt(mgrEnt)
	// if err != nil {
	// 	log.Errorf("gen gpatoken failed:%v", err)
	// 	return err
	// }

	gpaToken, claim := genJwtClaimFromMgrEnt(mgrEnt)
	jwtClaim := jwt.MapClaims(claim)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	tokenStr, err := jwtToken.SignedString([]byte(CookieKeySalt))
	if err != nil {
		log.Errorf("gen jwt token failed:%v", err)
		return err
	}

	// 设置cookie
	c.SetCookie(CookieKeyGpa, tokenStr, LoginCookieExpire, "", "", false, false)
	// c.SetCookie(CookieKeySign, sign, 3600*24*7, "", "", false, false)

	// 设置context
	gpaToken.MgrEnt = mgrEnt
	c.Set(CtxKeyGpa, gpaToken)

	// record login
	//go managermodel.RecordLoginStatus(mgrEnt)
	return nil
}

func ClearLogin(c *gin.Context) {
	c.SetCookie(CookieKeyGpa, "", -1, "", "", false, false)
	// c.SetCookie(CookieKeySign, "", -1, "", "", false, false)
}

// func ChangePwd(c *gin.Context) {
// 	c.HTML(http.StatusOK, "login/changepwd.tpl", gin.H{
// 		"LoginInfo": getLoginInfo(c),
// 		"IsLogin":   isLoginIn(c),
// 	})
// }

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
		JsonErr(c, ErrInvalidParam)
		return
	}

	loginEnt := getLoginInfo(c)
	if loginEnt == nil || loginEnt.Mid <= 0 {
		JsonErr(c, ErrNotLogin)
		return
	}

	mgrEnt, err := dao.GetManagerById(loginEnt.Mid)
	if err != nil {
		JsonErr(c, ErrMgrNotExist)
		return
	}

	if req.Email != mgrEnt.Email {
		JsonErrMsg(c, ErrInvalidParam, "请提交自己的邮箱")
		return
	}

	if !mgrEnt.IsPwdValid(req.Oldpwd) {
		JsonErr(c, ErrInvalidPass)
		return
	}

	mgrEnt.SetPwd(req.Newpwd)

	// _, err = base.UpdateByCol("mid", loginEnt.Mid, mgrDbEnt, []string{"mgrPwd"})
	err = mgrEnt.UpdateById([]string{"pwd"})

	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	JsonOk(c, gin.H{
		"id": mgrEnt.Id,
	})
	// c.Redirect(http.StatusMovedPermanently, RouteManagerList)
}

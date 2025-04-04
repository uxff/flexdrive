package apihandler

import (
	"time"

	"github.com/gin-gonic/gin"
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
	Username string `form:"username" json:"username"`
	Pwd      string `form:"password" json:"password"`
	Captcha  string `form:"captcha" json:"captcha"`
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
	mgrEnt, err := dao.GetManagerByEmail(req.Username)
	if err != nil {
		log.Errorf("query by email:%s failed:%v", req.Username, err)
		JsonErr(c, ErrMgrNotExist)
		return
	}

	if mgrEnt == nil || mgrEnt.Id == 0 {
		log.Warnf("email:%s not exist, verify failed", req.Username)
		JsonErr(c, ErrMgrNotExist)
		return
	}

	// 验证密码
	if !mgrEnt.IsPwdValid(req.Pwd) {
		log.Warnf("mgr pwd not matched, verify failed. email:%s", req.Username)
		JsonErr(c, ErrInvalidPass)
		return
	}

	// 判断账号是否已被禁用
	if mgrEnt.Status != base.StatusNormal {
		JsonErr(c, ErrMgrDisabled)
		return
	}

	// 登录成功 种下cookie
	token, _ := AcceptLogin(c, mgrEnt)

	JsonOk(c, gin.H{
		"date":      time.Now(),
		"API-Token": token,
		// "mgr":       mgrEnt,
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
	// 有问题，清理Cookie后，其实已经发放出去的Token还是有效的，可以继续访问
	ClearLogin(c)
	//StdResponse(c, ErrSuccess, nil)
	JsonOk(c, gin.H{
		"date": time.Now(),
	})
	// c.Redirect(http.StatusMovedPermanently, RouteLogin)
}

// 受理登录
func AcceptLogin(c *gin.Context, mgrEnt *dao.Manager) (token string, err error) {
	mgrEnt.LastLoginIp = c.Request.Header.Get("X-Real-IP")
	mgrEnt.LastLoginAt = time.Now() //util.JsonTime(time.Now()) // time.Now().Format("2006-01-02 15:04:05")

	// 生成签名
	// gpaToken, _, _, err := genGpaFromMgrEnt(mgrEnt)
	// if err != nil {
	// 	log.Errorf("gen gpatoken failed:%v", err)
	// 	return err
	// }

	gpaToken, jwtTokenStr, err := genJwtClaimFromMgrEnt(mgrEnt)
	// jwtClaim :=
	// jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claim))
	// tokenStr, err := jwtToken.SignedString([]byte(CookieKeySalt))
	if err != nil {
		log.Errorf("gen jwt token failed:%v", err)
		return "", err
	}

	// 设置cookie
	// c.SetCookie(CookieKeyGpa, jwtTokenStr, LoginCookieExpire, "", "", false, false) // use API-Token instead
	c.Header("API-Token", jwtTokenStr)
	c.SetCookie(CookieKeyGpa, jwtTokenStr, 3600*24*7, "", "", false, false)

	// 设置context
	gpaToken.MgrEnt = mgrEnt
	gpaToken.MgrEnt.Pwd = ""
	c.Set(CtxKeyGpa, gpaToken)

	// record login
	//go managermodel.RecordLoginStatus(mgrEnt)
	return jwtTokenStr, err
}

func ClearLogin(c *gin.Context) {
	// will be useless when using API-Token header
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

package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

type SignupRequest struct {
	Email   string `form:"email"`
	Pwd     string `form:"pwd"`
	RePwd   string `form:"repwd"`
	Captcha string `form:"captcha"`
}

func (r *SignupRequest) ToEnt() *dao.User {
	e := &dao.User{
		Email: r.Email,
		//Pwd: r.Pwd,
	}

	//e.SetPwd(r.Pwd)

	return e
}

func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "login/signup.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
	})
}

// 提交登录的处理
func SignupForm(c *gin.Context) {

	// 参数是否正确
	req := &SignupRequest{}
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

	log.Debugf("signup request:%+v", req)

	if req.Email == "" {
		JsonErrMsg(c, ErrInvalidParam, "邮箱不能为空")
		return
	}

	existEnt, err := dao.GetUserByEmail(req.Email)
	if err != nil {
		log.Errorf("query by email:%s failed:%v", req.Email, err)
		JsonErr(c, ErrInternal)
		return
	}

	if existEnt != nil && existEnt.Id > 0 {
		log.Warnf("email:%s already exist, cannot register again", req.Email)
		JsonErr(c, ErrEmailDuplicate)
		return
	}

	if len(req.Pwd) < 6 {
		JsonErrMsg(c, ErrInvalidPass, "密码太短")
		return
	}

	if req.Pwd != req.RePwd {
		JsonErrMsg(c, ErrInvalidPass, "2次密码不一致")
		return
	}

	userEnt := req.ToEnt()

	userEnt.SetPwd(req.Pwd)

	initialLevelEnt, err := dao.GetDefaultUserLevel()
	if err != nil || initialLevelEnt == nil || initialLevelEnt.Id <= 0 {
		log.Errorf("get default userlevel failed:%v", err)
		JsonErrMsg(c, ErrInternal, "没有默认等级，请联系管理员创建会员等级")
		return
	}

	// 初始等级及空间
	userEnt.LevelId = initialLevelEnt.Id
	userEnt.QuotaSpace = initialLevelEnt.QuotaSpace
	userEnt.Status = base.StatusNormal
	userEnt.LastLoginAt = time.Now()
	userEnt.LastLoginIp = getRemoteIp(c)

	_, err = base.Insert(userEnt)
	if err != nil {
		log.Errorf("create user failed:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	// 注册成功 种下cookie
	tokenStr, _ := AcceptLogin(c, userEnt)

	// c.Redirect(http.StatusMovedPermanently, RouteHome)
	//StdResponse(c, ErrSuccess, "/")
	JsonOk(c, gin.H{
		"API-Token": tokenStr,
	})
}

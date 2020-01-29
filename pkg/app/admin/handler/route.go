package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/uxff/flexdrive/pkg/log"

	"github.com/gin-gonic/gin"
)

const (
	RouteHome          = "/"
	RouteLogin         = "/login"
	RouteLogout        = "/logout"
	RouteManagerList   = "/manager/list"
	RouteRoleList      = "/role/list"
	RouteUserList      = "/user/list"
	RouteFileIndexList = "/file/list"
	RouteNodeList      = "/node/list"
	RouteShareList     = "/share/list"
	RouteUserLevelList = "/userlevel/list"
	RouteChangePwd     = "/changepwd"
)

var adminServer *http.Server
var router = gin.New() // *gin.Engine // 在本包init函数之前运行
var tplFuncMap = make(template.FuncMap, 0)

func init() {
	loadFuncMap()
}

func StartHttpServer(addr string) error {
	gin.SetMode(gin.DebugMode)

	hostName, _ := os.Hostname()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"hostname": hostName,
		})
	})

	// 公共路由
	// 登录
	router.GET("/login", TraceMiddleWare, Login)
	router.POST("/login", TraceMiddleWare, LoginForm)
	router.GET("/logout", TraceMiddleWare, Logout)
	//router.GET("/app/config", TraceMiddleWare, GetAppConfig)

	// 验证码
	router.GET("/captcha", GetCaptcha)

	// 导出下载 基于登录cookie验证
	authRouter := router.Group("/", TraceMiddleWare, AuthMiddleWare)
	authRouter.POST("/changePwd", ManagerChangePwd) // 修改自己的密码 不受角色限制
	authRouter.GET("/", Index)

	// 基础基于登录cookie并rabc授权的验证
	// 如果增加接口，必须在现有的菜单下，否则会被权限控制拦住
	// 也就是增加的接口必须以下面的group中的某一个路径开头
	rbacRouter := router.Group("/", TraceMiddleWare, AuthMiddleWare, RbacAuthMiddleWare)

	rbacRouter.GET("/role/add", RoleAdd)
	rbacRouter.POST("/role/add", RoleAddForm)
	rbacRouter.GET("/role/edit/:id", RoleEdit)
	rbacRouter.POST("/role/edit/:id", RoleAddForm)
	rbacRouter.GET("/role/enable/:id/:enable", RoleEnable)
	rbacRouter.GET("/role/list", RoleList)

	//rbacRouter.POST("/role/rbac/edit/:roleid", RoleRbacSet)
	//rbacRouter.GET("/role/rbac/list/:roleid", RoleRbacGet)

	rbacRouter.GET("/manager/list", ManagerList)
	rbacRouter.GET("/manager/add", ManagerAdd)
	rbacRouter.POST("/manager/add", ManagerAddForm)
	rbacRouter.GET("/manager/edit/:mid", ManagerEdit)
	rbacRouter.POST("/manager/edit/:mid", ManagerAddForm)
	//authRouter.POST("/manager/modifyPwd", ManagerChangePwd)
	rbacRouter.GET("/manager/enable/:mid/:enable", ManagerEnable)

	rbacRouter.GET("/user/list", UserList)
	rbacRouter.GET("/user/enable/:id/:enable", UserEnable)

	rbacRouter.GET("/userlevel/add", UserLevelAdd)
	rbacRouter.POST("/userlevel/add", UserLevelAddForm)
	rbacRouter.GET("/userlevel/edit/:id", UserLevelEdit)
	rbacRouter.POST("/userlevel/edit/:id", UserLevelAddForm)
	rbacRouter.GET("/userlevel/enable/:id/:enable", UserLevelEnable)
	rbacRouter.GET("/userlevel/list", UserLevelList)

	rbacRouter.GET("/share/list", ShareList)
	rbacRouter.GET("/share/enable/:id/:enable", ShareEnable)

	rbacRouter.GET("/node/list", UserList)
	// rbacRouter.GET("/share/enable/:id/:enable", UserEnable)

	rbacRouter.GET("/fileindex/list", FileIndexList)
	rbacRouter.GET("/fileindex/enable/:id/:enable", FileIndexEnable)

	adminServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	//fnMap := template.FuncMap{}
	//fnMap[""] =

	router.SetFuncMap(tplFuncMap)

	// gin的debug 模式下每次访问请求都会读取模板 release模式下不会
	router.LoadHTMLGlob("pkg/app/admin/view/**/*")

	// js 静态资源 在nginx下应该由nginx来服务比较专业
	router.StaticFS("/static", http.Dir("static"))

	return adminServer.ListenAndServe()
}

func ShutdownHttpServer() {
	if adminServer != nil {
		err := adminServer.Shutdown(context.Background())
		if err != nil {
			log.Errorf("shutdown http server failed:%v", err)
		}

		adminServer = nil
	}
}

func loadFuncMap() {
	tplFuncMap["i18nja"] = func(format string, args ...interface{}) string {
		return "" //i18n.Tr("ja-JP", format, args...)
	}
	//"i18n": i18n.Tr,
	tplFuncMap["datenow"] = func(format string) string {
		return time.Now().Add(time.Duration(9) * time.Hour).Format(format)
	}
	tplFuncMap["dateformatJst"] = func(in time.Time) string {
		in = in.Add(time.Duration(9) * time.Hour)
		return in.Format("2006/01/02 15:04")
	}

	tplFuncMap["qescape"] = func(in string) string {
		return url.QueryEscape(in)
	}
	tplFuncMap["nl2br"] = func(in string) string {
		return strings.Replace(in, "\n", "<br>", -1)
	}

	tplFuncMap["tostr"] = func(in interface{}) string {
		return fmt.Sprintf("%d", in) //convert.ToStr(reflect.ValueOf(in).Interface())
	}

	tplFuncMap["first"] = func(in interface{}) interface{} {
		return reflect.ValueOf(in).Index(0).Interface()
	}

	tplFuncMap["last"] = func(in interface{}) interface{} {
		s := reflect.ValueOf(in)
		return s.Index(s.Len() - 1).Interface()
	}

	tplFuncMap["truncate"] = func(in string, length int) string {
		return runewidth.Truncate(in, length, "...")
	}

	tplFuncMap["noname"] = func(in string) string {
		if in == "" {
			return "(未入力)"
		}
		return in
	}

	tplFuncMap["cleanurl"] = func(in string) string {
		return strings.Trim(strings.Trim(in, " "), "　")
	}

	tplFuncMap["append"] = func(data map[interface{}]interface{}, key string, value interface{}) template.JS {
		if _, ok := data[key].([]interface{}); !ok {
			data[key] = []interface{}{value}
		} else {
			data[key] = append(data[key].([]interface{}), value)
		}
		return template.JS("")
	}

	tplFuncMap["appendmap"] = func(data map[interface{}]interface{}, key string, name string, value interface{}) template.JS {
		v := map[string]interface{}{name: value}

		if _, ok := data[key].([]interface{}); !ok {
			data[key] = []interface{}{v}
		} else {
			data[key] = append(data[key].([]interface{}), v)
		}
		return template.JS("")
	}
	tplFuncMap["urlfor"] = func(endpoint string, values ...interface{}) string {
		return endpoint
	}
	tplFuncMap["captchaUrl"] = func() string {
		return fmt.Sprintf("/captcha?t=%d", time.Now().Unix())
	}
}

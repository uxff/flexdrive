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

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-runewidth"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

const (
	RouteHome         = "/"
	RouteLogin        = "/login"
	RouteLogout       = "/logout"
	RouteUserList     = "/user/list"
	RouteUserFileList = "/my/file/list"
	RouteShareList    = "/share/list"
	RouteChangePwd    = "/changePwd"
)

var customerServer *http.Server
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
	router.GET("/signup", TraceMiddleWare, Signup)
	router.POST("/signup", TraceMiddleWare, SignupForm)
	router.GET("/logout", TraceMiddleWare, Logout)
	router.GET("/share/search", TraceMiddleWare, ShareSearch)
	//router.GET("/app/config", TraceMiddleWare, GetAppConfig)
	router.GET("/", TraceMiddleWare, Index)

	// 验证码
	router.GET("/captcha", GetCaptcha)

	// 导出下载 基于登录cookie验证
	authRouter := router.Group("/", TraceMiddleWare, AuthMiddleWare)
	authRouter.GET("/changePwd", ChangePwd)      // 修改自己的密码 不受角色限制
	authRouter.POST("/changePwd", ChangePwdForm) // 修改自己的密码 不受角色限制

	// authRouter.GET("/user/list", UserList)
	// authRouter.GET("/user/enable/:id/:enable", UserEnable)

	authRouter.GET("/my/share/list", ShareList)
	authRouter.GET("/my/share/check/:userFileId", ShareCheck)
	authRouter.POST("/my/share/add", ShareAdd)
	authRouter.GET("/my/share/enable/:id/:enable", ShareEnable)

	authRouter.GET("/my/order/create", OrderCreate)
	authRouter.POST("/my/order/create", OrderCreateForm)
	authRouter.GET("/my/order/list", OrderList)
	authRouter.GET("/my/share/check/:userFileId", ShareCheck)
	authRouter.POST("/my/share/add", ShareAdd)
	authRouter.GET("/my/share/enable/:id/:enable", ShareEnable)

	authRouter.GET("/my/file/list", UserFileList)
	authRouter.POST("/my/file/newfolder", UserFileNewFolder)
	authRouter.POST("/my/file/upload", UploadForm)
	authRouter.GET("/my/file/enable/:id/:enable", UserFileEnable)

	customerServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	router.SetFuncMap(tplFuncMap)

	// gin的debug 模式下每次访问请求都会读取模板 release模式下不会
	router.LoadHTMLGlob("pkg/app/customer/view/**/*")

	// js 静态资源 在nginx下应该由nginx来服务比较专业
	router.StaticFS("/static", http.Dir("static"))

	return customerServer.ListenAndServe()
}

func ShutdownHttpServer() {
	if customerServer != nil {
		err := customerServer.Shutdown(context.Background())
		if err != nil {
			log.Errorf("shutdown http server failed:%v", err)
		}

		customerServer = nil
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
	tplFuncMap["mgrStatus"] = func(status int) string {
		return base.StatusMap[status]
	}
	// 所有的空间单位必须是int64
	tplFuncMap["space4Human"] = func(space int64) string {
		if space < 1024 {
			return fmt.Sprintf("%d kB", space)
		}
		if space < 1024*1024 {
			return fmt.Sprintf("%d MB", space/1024)
		}
		if space < 1024*1024*1024 {
			return fmt.Sprintf("%d GB", space/1024/1024)
		}
		return fmt.Sprintf("%d TB", space/1024/1024/1024)
	}
	// 所有的空间单位必须是int64
	tplFuncMap["size4Human"] = func(space int64) string {
		if space < 1024 {
			return fmt.Sprintf("%d B", space)
		}
		if space < 1024*1024 {
			return fmt.Sprintf("%.01f kB", float32(space)/1024)
		}
		if space < 1024*1024*1024 {
			return fmt.Sprintf("%.01f MB", float32(space)/1024/1024)
		}
		return fmt.Sprintf("%.02f GB", float32(space)/1024/1024/1024)
	}
	// 所有的空间单位必须是int64
	tplFuncMap["spaceRate"] = func(used int64, quota int64) string {
		return fmt.Sprintf("%d", int(float32(used)/float32(quota)*100))
	}
	tplFuncMap["orderStatus"] = func(orderStatus int) string {
		return dao.OrderStatusMap[orderStatus]
	}
}

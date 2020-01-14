package handler

import (
	"context"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/uxff/flexdrive/pkg/log"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var adminServer *http.Server
var router = gin.New() // *gin.Engine // 在本包init函数之前运行

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
	router.GET("/api/logout", TraceMiddleWare, Logout)
	//router.GET("/api/app/config", TraceMiddleWare, GetAppConfig)

	// 验证码
	router.GET("/api/captcha", GetCaptcha)

	// 导出下载 基于登录cookie验证
	authRouter := router.Group("/api", TraceMiddleWare, AuthMiddleWare)
	authRouter.POST("/manager/modifyPwd", ManagerChangePwd)

	// 基础基于登录cookie并rabc授权的验证
	// 如果增加接口，必须在现有的菜单下，否则会被权限控制拦住
	// 也就是增加的接口必须以下面的group中的某一个路径开头
	rbacRouter := router.Group("/api", TraceMiddleWare, AuthMiddleWare, RbacAuthMiddleWare)

	//rbacRouter.POST("/role/add", RoleAdd)
	//// rbacRouter.POST("/role/edit/:roleid", RoleAdd)
	//rbacRouter.POST("/role/delete/:roleid", RoleDelete)
	//rbacRouter.GET("/role/list", RoleList)
	//
	//rbacRouter.POST("/role/rbac/edit/:roleid", RoleRbacSet)
	//rbacRouter.GET("/role/rbac/list/:roleid", RoleRbacGet)

	rbacRouter.GET("/manager/list", ManagerList)
	rbacRouter.POST("/manager/add", ManagerAdd)
	// rbacRouter.POST("/manager/edit/:mid", ManagerAdd)
	rbacRouter.POST("/manager/enable/:mid/:enable", ManagerEnable)

	//rbacRouter.GET("/merchant/list", MerchantList)
	//rbacRouter.POST("/merchant/add", MerchantAdd)
	//// rbacRouter.POST("/merchant/edit/:merid", MerchantAdd)
	//rbacRouter.POST("/merchant/enable/:merid/:enable", MerchantEnable)
	//rbacRouter.GET("/merchant/export", MerchantExport)
	//
	//rbacRouter.GET("/agent/list", AgentList)
	//rbacRouter.POST("/agent/add", AgentAdd)
	//// rbacRouter.POST("/agent/edit/:agentid", AgentAdd)
	//rbacRouter.POST("/agent/delete/:agentid", AgentDelete)
	//
	//rbacRouter.GET("/bizaccount/list", BizAccountList)
	//rbacRouter.POST("/bizaccount/add", BizAccountAdd)
	//// rbacRouter.POST("/bizaccount/edit/:baid", BizAccountAdd)
	//rbacRouter.POST("/bizaccount/enable/:baid/:enable", BizAccountEnable)
	//rbacRouter.GET("/bizaccount/thirdchannels/list", ThirdChannelsList)

	adminServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	//fnMap := template.FuncMap{}
	//fnMap[""] =

	router.SetFuncMap(template.FuncMap{
		"i18nja": func(format string, args ...interface{}) string {
			return "" //i18n.Tr("ja-JP", format, args...)
		},
		//"i18n": i18n.Tr,
		"datenow": func(format string) string {
			return time.Now().Add(time.Duration(9) * time.Hour).Format(format)
		},
		"dateformatJst": func(in time.Time) string {
			in = in.Add(time.Duration(9) * time.Hour)
			return in.Format("2006/01/02 15:04")
		},

		"qescape": func(in string) string {
			return url.QueryEscape(in)
		},
		"nl2br": func(in string) string {
			return strings.Replace(in, "\n", "<br>", -1)
		},

		"tostr": func(in interface{}) string {
			return fmt.Sprintf("%d", in) //convert.ToStr(reflect.ValueOf(in).Interface())
		},

		"first": func(in interface{}) interface{} {
			return reflect.ValueOf(in).Index(0).Interface()
		},

		"last": func(in interface{}) interface{} {
			s := reflect.ValueOf(in)
			return s.Index(s.Len() - 1).Interface()
		},

		"truncate": func(in string, length int) string {
			return runewidth.Truncate(in, length, "...")
		},

		"noname": func(in string) string {
			if in == "" {
				return "(未入力)"
			}
			return in
		},

		"cleanurl": func(in string) string {
			return strings.Trim(strings.Trim(in, " "), "　")
		},

		"append": func(data map[interface{}]interface{}, key string, value interface{}) template.JS {
			if _, ok := data[key].([]interface{}); !ok {
				data[key] = []interface{}{value}
			} else {
				data[key] = append(data[key].([]interface{}), value)
			}
			return template.JS("")
		},

		"appendmap": func(data map[interface{}]interface{}, key string, name string, value interface{}) template.JS {
			v := map[string]interface{}{name: value}

			if _, ok := data[key].([]interface{}); !ok {
				data[key] = []interface{}{v}
			} else {
				data[key] = append(data[key].([]interface{}), v)
			}
			return template.JS("")
		},
		"urlfor": func(endpoint string, values ...interface{}) string {
			return endpoint
		},
	})

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

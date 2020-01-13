package handler

import (
	"context"
	"github.com/uxff/flexdrive/pkg/log"
	"net/http"
	"os"

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
	router.POST("/api/login", TraceMiddleWare, Login)
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

	rbacRouter.POST("/role/add", RoleAdd)
	// rbacRouter.POST("/role/edit/:roleid", RoleAdd)
	rbacRouter.POST("/role/delete/:roleid", RoleDelete)
	rbacRouter.GET("/role/list", RoleList)

	rbacRouter.POST("/role/rbac/edit/:roleid", RoleRbacSet)
	rbacRouter.GET("/role/rbac/list/:roleid", RoleRbacGet)

	rbacRouter.GET("/manager/list", ManagerList)
	rbacRouter.POST("/manager/add", ManagerAdd)
	// rbacRouter.POST("/manager/edit/:mid", ManagerAdd)
	rbacRouter.POST("/manager/enable/:mid/:enable", ManagerEnable)

	rbacRouter.GET("/merchant/list", MerchantList)
	rbacRouter.POST("/merchant/add", MerchantAdd)
	// rbacRouter.POST("/merchant/edit/:merid", MerchantAdd)
	rbacRouter.POST("/merchant/enable/:merid/:enable", MerchantEnable)
	rbacRouter.GET("/merchant/export", MerchantExport)

	rbacRouter.GET("/agent/list", AgentList)
	rbacRouter.POST("/agent/add", AgentAdd)
	// rbacRouter.POST("/agent/edit/:agentid", AgentAdd)
	rbacRouter.POST("/agent/delete/:agentid", AgentDelete)

	rbacRouter.GET("/bizaccount/list", BizAccountList)
	rbacRouter.POST("/bizaccount/add", BizAccountAdd)
	// rbacRouter.POST("/bizaccount/edit/:baid", BizAccountAdd)
	rbacRouter.POST("/bizaccount/enable/:baid/:enable", BizAccountEnable)
	rbacRouter.GET("/bizaccount/thirdchannels/list", ThirdChannelsList)

	adminServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// gin的debug 模式下每次访问请求都会读取模板 release模式下不会
	router.LoadHTMLGlob("view/*/*")

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

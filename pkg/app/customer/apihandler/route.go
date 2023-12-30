package handler

import (
	"github.com/gin-gonic/gin"
)

func LoadRouter(router *gin.RouterGroup) {
	//gin.SetMode(gin.DebugMode)

	// gin的debug 模式下每次访问请求都会读取模板 release模式下不会
	// rootRouter.LoadHTMLGlob("pkg/app/customer/view/**/*")

	// js 静态资源 在nginx下应该由nginx来服务比较专业
	// rootRouter.StaticFS("/static", http.Dir("static"))

	// var router = rootRouter.Group(assignedGroupPrefix)

	// hostName, _ := os.Hostname()
	// router.GET("/health", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"status":   "ok",
	// 		"hostname": hostName,
	// 	})
	// })

	// 公共路由
	// 登录
	// router.GET("/login", TraceMiddleWare, Login)
	router.POST("/login", TraceMiddleWare, LoginForm)
	router.GET("/signup", TraceMiddleWare, Signup)
	router.POST("/signup", TraceMiddleWare, SignupForm)
	router.GET("/logout", TraceMiddleWare, Logout)
	router.GET("/share/search", TraceMiddleWare, ShareSearch)
	router.GET("/s/:shareHash", TraceMiddleWare, ShareDetail)
	//router.GET("/app/config", TraceMiddleWare, GetAppConfig)
	router.GET("/", TraceMiddleWare, Index)
	router.GET("/file/:fileHash/:fileName", Fs)

	// 验证码
	router.GET("/captcha", GetCaptcha)

	// 导出下载 基于登录cookie验证
	authRouter := router.Group("/", TraceMiddleWare, AuthMiddleWare)
	// authRouter.GET("/changePwd", ChangePwd)      // 修改自己的密码 不受角色限制
	authRouter.POST("/changePwd", ChangePwdForm) // 修改自己的密码 不受角色限制

	// authRouter.GET("/user/list", UserList)
	// authRouter.GET("/user/enable/:id/:enable", UserEnable)

	authRouter.GET("/my/share/list", ShareList)
	authRouter.GET("/my/share/check/:userFileId", ShareCheck)
	authRouter.POST("/my/share/add", ShareAdd)
	authRouter.GET("/my/share/enable/:id/:enable", ShareEnable)

	// authRouter.GET("/my/order/create", OrderCreate)
	authRouter.POST("/my/order/create", OrderCreateForm)
	authRouter.GET("/my/order/detail/:orderId", OrderDetail)
	authRouter.GET("/my/order/mockpay/:orderId", Mockpay)
	authRouter.POST("/my/order/notify", MockpayForm)
	authRouter.GET("/my/order/list", OrderList)

	authRouter.GET("/my/offlinetask/list", OfflineTaskList)
	authRouter.GET("/my/offlinetask/add", OfflineTaskAdd)
	authRouter.GET("/my/offlinetask/enable/:id/:enable", OfflineTaskEnable)

	authRouter.GET("/my/profile", Profile)

	authRouter.GET("/my/file/list", UserFileList)
	authRouter.POST("/my/file/newfolder", UserFileNewFolder)
	authRouter.POST("/my/file/rename", UserFileRename)
	authRouter.POST("/my/file/upload", UploadForm)
	authRouter.GET("/my/file/enable/:id/:enable", UserFileEnable)

	// customerServer = &http.Server{
	// 	Addr:    addr,
	// 	Handler: router,
	// }

	// router.SetFuncMap(tplFuncMap)

	// return customerServer.ListenAndServe()
}

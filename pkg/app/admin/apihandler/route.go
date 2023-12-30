package apihandler

import (
	"github.com/gin-gonic/gin"
)

const (
	RouteHome          = "/"
	RouteLogin         = "/login"
	RouteLogout        = "/logout"
	RouteManagerList   = "/manager/list"
	RouteRoleList      = "/role/list"
	RouteUserList      = "/user/list"
	RouteUserFileList  = "/user/file/list"
	RouteFileIndexList = "/file/list"
	RouteNodeList      = "/node/list"
	RouteShareList     = "/share/list"
	RouteUserLevelList = "/userlevel/list"
	RouteOrderList     = "/order/list"
	RouteChangePwd     = "/changepwd"
)

// var adminServer *http.Server

// var router = gin.New() // *gin.Engine // 在本包init函数之前运行
// var tplFuncMap = make(template.FuncMap, 0)

func init() {
	// loadFuncMap()
}

func LoadRouter(rootRouter *gin.Engine, assignedGroupPrefix string) error {
	//gin.SetMode(gin.DebugMode)

	var router = rootRouter.Group(assignedGroupPrefix)

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
	// authRouter.GET("/changePwd", ChangePwd)      // 修改自己的密码 不受角色限制
	authRouter.POST("/changePwd", ChangePwdForm) // 修改自己的密码 不受角色限制
	authRouter.GET("/", Index)
	authRouter.GET("/file/:fileHash/:fileName", Fs)

	// 基础基于登录cookie并rabc授权的验证
	// 如果增加接口，必须在现有的菜单下，否则会被权限控制拦住
	// 也就是增加的接口必须以下面的group中的某一个路径开头
	rbacRouter := router.Group("/", TraceMiddleWare, AuthMiddleWare, RbacAuthMiddleWare)

	// rbacRouter.GET("/role/add", RoleAdd)
	rbacRouter.POST("/role/add", RoleAddForm)
	// rbacRouter.GET("/role/edit/:id", RoleEdit)
	rbacRouter.POST("/role/edit/:id", RoleAddForm)
	rbacRouter.POST("/role/enable/:id/:enable", RoleEnable)
	rbacRouter.GET("/role/list", RoleList)

	//rbacRouter.POST("/role/rbac/edit/:roleid", RoleRbacSet)
	//rbacRouter.GET("/role/rbac/list/:roleid", RoleRbacGet)

	rbacRouter.GET("/manager/list", ManagerList)
	// rbacRouter.GET("/manager/add", ManagerAdd)
	rbacRouter.POST("/manager/add", ManagerAddForm)
	// rbacRouter.GET("/manager/edit/:mid", ManagerEdit)
	rbacRouter.POST("/manager/edit/:mid", ManagerAddForm)
	//authRouter.POST("/manager/modifyPwd", ManagerChangePwd)
	rbacRouter.GET("/manager/enable/:mid/:enable", ManagerEnable)

	rbacRouter.GET("/user/list", UserList)
	rbacRouter.GET("/user/enable/:id/:enable", UserEnable)
	rbacRouter.GET("/user/file/list", UserFileList)
	rbacRouter.GET("/user/file/enable/:id/:enable", UserFileEnable)

	// rbacRouter.GET("/userlevel/add", UserLevelAdd)
	rbacRouter.POST("/userlevel/add", UserLevelAddForm)
	// rbacRouter.GET("/userlevel/edit/:id", UserLevelEdit)
	rbacRouter.POST("/userlevel/edit/:id", UserLevelAddForm)
	rbacRouter.GET("/userlevel/enable/:id/:enable", UserLevelEnable)
	rbacRouter.GET("/userlevel/list", UserLevelList)

	rbacRouter.GET("/share/list", ShareList)
	rbacRouter.GET("/share/enable/:id/:enable", ShareEnable)

	rbacRouter.GET("/order/list", OrderList)
	rbacRouter.GET("/order/refund/:id", OrderRefund)

	rbacRouter.GET("/node/list", NodeList)
	rbacRouter.POST("/node/setspace", NodeSetspace)

	rbacRouter.GET("/fileindex/list", FileIndexList)
	rbacRouter.GET("/fileindex/enable/:id/:enable", FileIndexEnable)

	// adminServer = &http.Server{
	// 	Addr:    addr,
	// 	Handler: router,
	// }

	// router.SetFuncMap(tplFuncMap)

	// // gin的debug 模式下每次访问请求都会读取模板 release模式下不会
	// router.LoadHTMLGlob("pkg/app/admin/view/**/*")

	// // js 静态资源 在nginx下应该由nginx来服务比较专业
	// router.StaticFS("/static", http.Dir("static"))

	// return adminServer.ListenAndServe()
	return nil
}

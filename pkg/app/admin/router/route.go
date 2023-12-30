package router

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	adminApiHadler "github.com/uxff/flexdrive/pkg/app/admin/apihandler"
	adminHadler "github.com/uxff/flexdrive/pkg/app/admin/handler"
	"github.com/uxff/flexdrive/pkg/log"
)

var adminServer *http.Server

func StartHttpServer(addr string) error {
	var router = gin.New()
	//gin.SetMode(gin.DebugMode)

	hostName, _ := os.Hostname()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"hostname": hostName,
		})
	})

	// 页面路由
	// pageRouter := router.Group("/", adminHadler.TraceMiddleWare)
	adminHadler.LoadRouter(router, "/")

	// API 路由
	// apiRouter := router.Group("/api", adminHadler.TraceMiddleWare)
	adminApiHadler.LoadRouter(router, "/api")

	adminServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// router.SetFuncMap(tplFuncMap)

	// // gin的debug 模式下每次访问请求都会读取模板 release模式下不会
	// router.LoadHTMLGlob("pkg/app/admin/view/**/*")

	// // js 静态资源 在nginx下应该由nginx来服务比较专业
	// router.StaticFS("/static", http.Dir("static"))

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
/*
*

		分布式(distributed)
	    运行方式gin web:
		SERVEADMIN=0.0.0.0:10011 SERVECUSTOMER=0.0.0.0:10012 SERVECLUSTER=0.0.0.0:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='mysql://user:pwd@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4&parseTime=True&loc=Local'  STORAGEDIR=./data/ ./main
		SERVEADMIN=0.0.0.0:10011 SERVECUSTOMER=0.0.0.0:10012 SERVECLUSTER=0.0.0.0:10013 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./main
	    运行方式gin api(no web, web should be served by react.js/vue.js):
		SERVEAPI=0.0.0.0:10011 SERVECLUSTER=0.0.0.0:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='mysql://user:pwd@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4&parseTime=True&loc=Local'  STORAGEDIR=./data/ ./main

		for cluster:
		SERVEWEB=0.0.0.0:10011 SERVECLUSTER=0.0.0.0:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./main
		SERVEWEB=0.0.0.0:10021 SERVECLUSTER=0.0.0.0:10023 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./main
		SERVEWEB=0.0.0.0:10031 SERVECLUSTER=0.0.0.0:10033 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./main
*/
package main

import (
	"flag"
	slog "log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/app/nodestorage/model/storagemodel"

	"github.com/uxff/flexdrive/pkg/common"

	adminapihandler "github.com/uxff/flexdrive/pkg/app/admin/apihandler"
	adminhandler "github.com/uxff/flexdrive/pkg/app/admin/handler"
	customerapihandler "github.com/uxff/flexdrive/pkg/app/customer/apihandler"
	customerhandler "github.com/uxff/flexdrive/pkg/app/customer/handler"
	"github.com/uxff/flexdrive/pkg/envinit"
	"github.com/uxff/flexdrive/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	version     = "0.1"
	showVersion bool
	logLevel    = -1
	// default values, you can set these with env
	serveApi       = "0.0.0.0:10010" // should be set
	serveAdmin     = ""              //"0.0.0.0:10011"
	serveCustomer  = ""              //"0.0.0.0:10012"
	serveCluster   = "0.0.0.0:10013" // must be set
	clusterMembers = ""              //"0.0.0.0:10013,0.0.0.0:10023,0.0.0.0:10033"
	clusterId      = "flexdrive"
	dataDsn        = "mysql://user:pass@tcp(0.0.0.0:3306)/flexdrive?charset=utf8mb4&parseTime=True&loc=Local"
	cacheDsn       = ""
	storageDir     = "/tmp/flexdrive/"
)

func main() {

	flag.IntVar(&logLevel, "l", logLevel, "log logLevel, -1:debug, 0:info, 1:warn, 2:error")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()

	if showVersion {
		flag.Usage()
		os.Exit(0)
	}

	lcf := zap.NewDevelopmentConfig()
	lcf.Level.SetLevel(zapcore.Level(logLevel))
	lcf.Development = true
	// lcf.DisableStacktrace = true
	// logger, err := lcf.Build(zap.AddCallerSkip(2), zap.AddCaller())
	logger, err := lcf.Build(zap.AddCaller())
	if err != nil {
		slog.Fatalln("new log err:", err.Error())
	}

	log.SetLogger(logger.Sugar())

	if s := os.Getenv("DATADSN"); s != "" {
		log.Debugf("the datadsn from env: %s", s)
		dataDsn = s
	}

	if s := os.Getenv("SERVEAPI"); s != "" {
		log.Debugf("the serverapi from env: %s", s)
		serveApi = s
	}

	if s := os.Getenv("SERVEADMIN"); s != "" {
		log.Debugf("the serveradmin from env: %s", s)
		serveAdmin = s
	}

	if s := os.Getenv("SERVECUSTOMER"); s != "" {
		log.Debugf("the servercustomer from env: %s", s)
		serveCustomer = s
	}

	if s := os.Getenv("STORAGEDIR"); s != "" {
		log.Debugf("the STORAGEDIR from env: %s", s)
		storageDir = s
	}

	if s := os.Getenv("CLUSTERMEMBERS"); s != "" {
		log.Debugf("the CLUSTERMEMBERS from env: %s", s)
		clusterMembers = s
	}

	if s := os.Getenv("CLUSTERID"); s != "" {
		log.Debugf("the storageDir from env: %s", s)
		clusterId = s
	}

	if s := os.Getenv("SERVECLUSTER"); s != "" {
		log.Debugf("the storageDir from env: %s", s)
		serveCluster = s
	}

	err = envinit.InitDb(common.DBMysqlDrive, dataDsn)
	if err != nil {
		log.Fatalf("cannot init db, err:%s", err)
	}

	log.Infof("db %s init ok", dataDsn)

	envMap := make(map[string]string)

	if err := Serve(envMap); err != nil {
		log.Fatalf("main error:%v", err)
	}
}

func Serve(envMap map[string]string) error {
	errCh := make(chan error, 1)

	wg := &sync.WaitGroup{}

	if serveAdmin != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- adminhandler.StartHttpServer(serveAdmin)
		}()
	}

	if serveCustomer != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- customerhandler.StartHttpServer(serveCustomer)
		}()
	}

	if serveApi != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- ServeApi(serveApi)
		}()
	}

	// todo handle when non cluster
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- storagemodel.StartNode(storageDir, serveCluster, clusterId, clusterMembers)
	}()

	select {
	case e := <-errCh:
		return e
	}

	wg.Wait()
	return nil
}

// ServeApi - serve api for admin and customer, will be used by react.js or vue.js
func ServeApi(addr string) error {
	router := gin.New()

	// customer
	customerApiRouter := router.Group("/api")

	// admin
	adminApiRouter := router.Group("/admapi")

	customerapihandler.LoadRouter(customerApiRouter)
	adminapihandler.LoadRouter(adminApiRouter)

	router.GET("/health", func(c *gin.Context) {
		hostName, _ := os.Hostname()
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"hostname": hostName,
			"info":     "this is flexdrive api server",
		})
	})

	webServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return webServer.ListenAndServe()
}

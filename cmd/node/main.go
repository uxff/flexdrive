/**
	分布式(distributed)
    运行方式：
 	APPENV=beta SERVEADMIN=127.0.0.1:10011 SERVECUSTOMER=127.0.0.1:10012 DATADSN='mysql://yourusername:yourpwd@tcp(yourmysqlhost)/yourdbname?charset=utf8mb4&parseTime=True&loc=Local' ./main
 	APPENV=beta SERVEADMIN=127.0.0.1:10011 SERVECUSTOMER=127.0.0.1:10012 DATADSN='sqlite3://./flexdrive.db' ./main
*/
package main

import (
	"flag"
	slog "log"
	"os"
	"sync"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/model/storagemodel"

	"github.com/uxff/flexdrive/pkg/common"

	adminhandler "github.com/uxff/flexdrive/pkg/app/admin/handler"
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
	serveCustomer = "127.0.0.1:10012"
	serveAdmin    = "127.0.0.1:10011"
	dataDsn       = "mysql://user:pass@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4&parseTime=True&loc=Local"
	cacheDsn      = ""
	storageDir    = "/tmp/flexdrive/"
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
	lcf.DisableStacktrace = true
	logger, err := lcf.Build(zap.AddCallerSkip(1))
	if err != nil {
		slog.Fatalln("new log err:", err.Error())
	}

	log.SetLogger(logger.Sugar())

	if s := os.Getenv("DATADSN"); s != "" {
		log.Debugf("the datadsn from env: %s", s)
		dataDsn = s
	}

	if s := os.Getenv("SERVEADMIN"); s != "" {
		log.Debugf("the serveradmin from env: %s", s)
		serveAdmin = s
	}

	if s := os.Getenv("SERVECUSTOMER"); s != "" {
		log.Debugf("the servercustomer from env: %s", s)
		serveCustomer = s
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
	go func() {
		wg.Add(1)
		defer wg.Done()
		errCh <- adminhandler.StartHttpServer(serveAdmin)
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		errCh <- customerhandler.StartHttpServer(serveCustomer)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		storagemodel.StartNode("me", storageDir)
	}()

	select {
	case e := <-errCh:
		return e
	}

	wg.Wait()
	return nil
}

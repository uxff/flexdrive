/**
	分布式(distributed)
    运行方式：
 	APPENV=beta SERVICE_PORT=9001 DATASTORE_URL='mysql://yourusername:yourpwd@tcp(yourmysqlhost)/yourdbname?charset=utf8mb4&parseTime=True&loc=Local' ./main
*/
package main

import (
	"flag"
	"github.com/uxff/flexdrive/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	slog "log"
	"os"
)

var (
	version         = "0.1"
	showVersion 	bool
	logLevel     	= 0
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
	lcf.Development = false
	lcf.DisableStacktrace = true
	logger, err := lcf.Build(zap.AddCallerSkip(1))
	if err != nil {
		slog.Fatalln("new log err:", err.Error())
	}

	log.SetLogger(logger.Sugar())

	envMap := make(map[string]string)

	if err := Serve(envMap); err != nil {
		log.Fatalf("main error:%v", err)
	}
}


func Serve(envMap map[string]string) error {

}

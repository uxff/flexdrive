package envinit

import (
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/uxff/flexdrive/pkg/log"
	"regexp"
	"time"
)

// Dbs 分库分表数据库名称映射
var Dbs map[string]*xorm.Engine

func init() {
	Dbs = make(map[string]*xorm.Engine)
}

// InitMysql 链接数据库 path 为 dsn
func InitMysql(path string) *xorm.Engine {
	var err error
	engine, err := xorm.NewEngine("mysql", path)
	if err != nil {

		log.Fatalf("xorm create err:", err)
		return nil
	}

	engine.Ping()
	engine.SetMaxIdleConns(20)
	engine.SetConnMaxLifetime(9 * time.Second)

	re := regexp.MustCompile(`/\w*\?`)
	str := re.FindString(path)
	if len(str) < 2 {
		log.Fatalf("prefix parse dbname err:", str, path)
		return nil
	}
	dbPrefix := ""
	mptable := core.NewPrefixMapper(&core.SnakeMapper{}, dbPrefix)
	engine.SetTableMapper(mptable)

	// Engine.ShowSQL(true)
	return engine
}

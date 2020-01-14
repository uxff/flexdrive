package envinit

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/uxff/flexdrive/pkg/log"

	_ "github.com/go-sql-driver/mysql"
)

// Dbs 分库分表数据库名称映射
var Dbs map[string]*xorm.Engine

func init() {
	Dbs = make(map[string]*xorm.Engine)
}

func InitMysql(namespace, dsn string) error {
	eng, err := ConnectMysql(dsn)
	if err != nil {
		log.Errorf("connect mysql %s error:%v", dsn, err)
		return err
	}

	// redo register namespace is not allowed
	if _, ok := Dbs[namespace]; ok {
		log.Errorf("namespace already exist, do not redo this")
		return errors.New("namespace already exit")
	}

	Dbs[namespace] = eng
	log.Debugf("namespace %s is registered ok", namespace)
	return nil
}

// InitMysql 链接数据库 path 为 dsn
func ConnectMysql(path string) (*xorm.Engine, error) {
	var err error
	engine, err := xorm.NewEngine("mysql", path)
	if err != nil {
		//log.Fatalf("xorm create err:", err)
		return nil, err
	}

	engine.Ping()
	engine.SetMaxIdleConns(20)
	engine.SetConnMaxLifetime(9 * time.Second)

	re := regexp.MustCompile(`/\w*\?`)
	str := re.FindString(path)
	if len(str) < 2 {
		//log.Fatalf("prefix parse dbname err:", str, path)
		return nil, errors.New("regexp find dbname failed")
	}
	dbPrefix := ""
	mptable := core.NewPrefixMapper(&core.SnakeMapper{}, dbPrefix)
	engine.SetTableMapper(mptable)

	// Engine.ShowSQL(true)
	return engine, nil
}

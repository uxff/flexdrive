package envinit

import (
	"errors"
	//"regexp"
	"strings"
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

// dsn 带 mysql://
func InitDb(namespace, dsn string) error {
	log.Debugf("will connect %s", dsn)
	// redo register namespace is not allowed
	if _, ok := Dbs[namespace]; ok {
		log.Errorf("namespace already exist, do not redo this")
		return errors.New("namespace already exit")
	}

	dsnPaths := strings.Split(dsn, "://")

	if len(dsnPaths) != 2 {
		return errors.New("dsn path must be like: mysql://user:pwd@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4")
	}
	engineType := dsnPaths[0]
	dsnPath := dsnPaths[1]

	var eng *xorm.Engine
	var err error
	switch engineType {
	case "mysql":
		eng, err = ConnectByEngine(engineType, dsnPath)
		if err != nil {
			log.Errorf("connect db %s error:%v", dsn, err)
			return err
		}
		if eng != nil {
			//eng.Exec("set sessoin sql_mode='NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';")
		}
	case "sqlite3":
		eng, err = ConnectByEngine(engineType, dsnPath)
		if err != nil {
			log.Errorf("connect db %s error:%v", dsn, err)
			return err
		}
	}

	if eng == nil {
		log.Errorf("connect db %s error:%v", dsn, err)
		return err
	}

	// eng, err := ConnectMysql(dsnPath)

	Dbs[namespace] = eng
	log.Debugf("namespace %s is registered ok", namespace)
	return nil
}

// InitMysql 链接数据库 path 为 dsn 带mysql://
/**
* @param en 引擎类型 mysql 或 sqlite3
 */
func ConnectByEngine(en string, dsnPath string) (*xorm.Engine, error) {
	var err error

	engine, err := xorm.NewEngine(en, dsnPath)
	if err != nil {
		//log.Fatalf("xorm create err:", err)
		return nil, err
	}

	engine.Ping()
	engine.SetMaxIdleConns(20)
	engine.SetConnMaxLifetime(9 * time.Second)

	// re := regexp.MustCompile(`/\w+\?`) //regexp.MustCompile(`/\w*\?`)
	// str := re.FindString(dsnPath)
	// if len(str) < 2 {
	// 	//log.Fatalf("prefix parse dbname err:", str, dsnPath)
	// 	return nil, errors.New("regexp find dbname failed")
	// }
	dbPrefix := ""
	mptable := core.NewPrefixMapper(&core.SnakeMapper{}, dbPrefix)
	engine.SetTableMapper(mptable)
	engine.SetColumnMapper(&core.SameMapper{})

	engine.ShowSQL(true)
	engine.ShowExecTime(true)
	// engine.Logger().SetLevel(core.LOG_DEBUG)
	return engine, nil
}

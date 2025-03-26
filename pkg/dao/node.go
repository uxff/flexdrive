package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/common"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/envinit"
	"github.com/uxff/flexdrive/pkg/log"
)

type Node struct {
	Id             int       `xorm:"not null pk autoincr comment('主键id') INT(10)"`
	NodeName       string    `xorm:"not null comment('节点名') VARCHAR(16)"`
	NodeAddr       string    `xorm:"not null comment('节点在集群中的通信地址') VARCHAR(32)"`
	ClusterId      string    `xorm:"not null comment('集群id') VARCHAR(40)"`
	Follow         string    `xorm:"not null comment('node who I follow') VARCHAR(40)"`
	TotalSpace     int64     `xorm:"not null default 0 comment('全部空间 单位KB') BIGINT(20)"`
	UsedSpace      int64     `xorm:"not null default 0 comment('使用的空间 单位KB') BIGINT(20)"`
	UnusedSpace    int64     `xorm:"not null default 0 comment('未使用的空间 单位KB') BIGINT(20)"`
	FileCount      int64     `xorm:"not null default 0 comment('文件数量') BIGINT(20)"`
	Remark         string    `xorm:"not null comment('房间备注') TEXT"`
	Created        time.Time `xorm:"created not null default '1999-12-31 00:00:00' comment('添加时间') TIMESTAMP"`
	Updated        time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status         int       `xorm:"not null default 1 comment('状态 1=正常启用 2=注册超时 99=删除 ') TINYINT(4)"`
	LastRegistered time.Time `xorm:"not null default '1999-12-31 00:00:00' comment('最后注册时间') TIMESTAMP"` // 声明的default无用 实际效果是CURRENT_TIMESTAMP
	// 上面时间xorm申明必须default null 数据库ddl可以不是default null, 因为 zero值问题
}

func (t Node) TableName() string {
	return "node"
}

func (t *Node) GetById(id int) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

func (t *Node) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func GetNodeById(id int) (*Node, error) {
	e := &Node{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func GetNodeByWorkerId(id string) (*Node, error) {
	e := &Node{}
	exist, err := base.GetByCol("nodeName", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

// 统计使用空间
func (t *Node) SumSpace() int64 {
	dbname := common.DBMysqlDrive // t.DbNamespace()
	// total, err := engine.Where("id >?", 1).SumInt(ss, "money")
	total1, err := envinit.Dbs[dbname].Where("nodeId=? and status=?", t.Id, 1).SumInt(&FileIndex{}, "space")
	if err != nil {
		log.Errorf("sum node(%d).usedSpace failed:%v", t.Id, err)
		return 0
	}
	total2, err := envinit.Dbs[dbname].Where("nodeId2=? and status=?", t.Id, 1).SumInt(&FileIndex{}, "space")
	if err != nil {
		log.Errorf("sum node(%d).usedSpace failed:%v", t.Id, err)
		return 0
	}
	total3, err := envinit.Dbs[dbname].Where("nodeId3=? and status=?", t.Id, 1).SumInt(&FileIndex{}, "space")
	if err != nil {
		log.Errorf("sum node(%d).usedSpace failed:%v", t.Id, err)
		return 0
	}
	t.UsedSpace = total1 + total2 + total3
	t.UnusedSpace = t.TotalSpace - t.UsedSpace
	err = t.UpdateById([]string{"usedSpace", "unusedSpace"})
	if err != nil {
		log.Errorf("sum node(%d).usedSpace failed:%v", t.Id, err)
		return 0
	}

	return t.UsedSpace
}

// 统计文件数
func (t *Node) CountFiles() int64 {
	dbname := common.DBMysqlDrive // t.DbNamespace()
	// total, err := engine.Where("id >?", 1).SumInt(ss, "money")
	total1, err := envinit.Dbs[dbname].Where("nodeId=? and status=?", t.Id, 1).Count(&FileIndex{}, "id")
	if err != nil {
		log.Errorf("count node(%d).files failed:%v", t.Id, err)
		return 0
	}
	total2, err := envinit.Dbs[dbname].Where("nodeId2=? and status=?", t.Id, 1).Count(&FileIndex{}, "id")
	if err != nil {
		log.Errorf("count node(%d).files failed:%v", t.Id, err)
		return 0
	}
	total3, err := envinit.Dbs[dbname].Where("nodeId3=? and status=?", t.Id, 1).Count(&FileIndex{}, "id")
	if err != nil {
		log.Errorf("count node(%d).files failed:%v", t.Id, err)
		return 0
	}
	t.FileCount = total1 + total2 + total3
	err = t.UpdateById([]string{"fileCount"})
	if err != nil {
		log.Errorf("count node(%d).files failed:%v", t.Id, err)
		return 0
	}

	return t.FileCount
}

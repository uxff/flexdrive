package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type Node struct {
	Id             int       `xorm:"not null pk autoincr comment('主键id') INT(10)"`
	NodeName       string    `xorm:"not null comment('节点名') VARCHAR(16)"`
	NodeAddr       string    `xorm:"not null comment('节点在集群中的通信地址') VARCHAR(32)"`
	ClusterId      string    `xorm:"not null comment('集群id') VARCHAR(40)"`
	TotalSpace     int64     `xorm:"not null default 0 comment('全部空间 单位KB') BIGINT(20)"`
	UsedSpace      int64     `xorm:"not null default 0 comment('使用的空间 单位KB') BIGINT(20)"`
	UnusedSpace    int64     `xorm:"not null default 0 comment('未使用的空间 单位KB') BIGINT(20)"`
	FileCount      int64     `xorm:"not null default 0 comment('文件数量') BIGINT(20)"`
	Remark         string    `xorm:"not null comment('房间备注') TEXT"`
	Created        time.Time `xorm:"created not null default '0000-00-00 00:00:00' comment('添加时间') TIMESTAMP"`
	Updated        time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status         int       `xorm:"not null default 1 comment('状态 1=正常启用 2=注册超时 99=删除 ') TINYINT(4)"`
	LastRegistered time.Time `xorm:"not null default ('0000-00-00 00:00:00') comment('最后注册时间') TIMESTAMP"` // 声明的default无用 实际效果是CURRENT_TIMESTAMP
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

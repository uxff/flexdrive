package flexdrivecode

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type Node struct {
	Id             int       `xorm:"not null pk autoincr comment('主键id') INT(10)"`
	NodeName       string    `xorm:"not null comment('节点号') VARCHAR(16)"`
	TotalSpace     int64     `xorm:"not null default 0 comment('全部空间 单位KB') BIGINT(20)"`
	UsedSpace      int64     `xorm:"not null default 0 comment('使用的空间 单位KB') BIGINT(20)"`
	FileCount      int64     `xorm:"not null default 0 comment('文件数量') BIGINT(20)"`
	Remark         string    `xorm:"not null comment('房间备注') TEXT"`
	Created        time.Time `xorm:"not null default '1999-12-31 00:00:00' comment('添加时间') TIMESTAMP"`
	Updated        time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status         int       `xorm:"not null default 1 comment('状态 1=正常启用 2=注册超时 99=删除 ') TINYINT(4)"`
	LastRegistered time.Time `xorm:"not null default '1999-12-31 00:00:00' comment('最后注册时间') TIMESTAMP"`
}

func (t Node) TableName() string {
	return "node"
}

func (t *Node) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

package flexdrivecode

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type OperLog struct {
	Id          int       `xorm:"not null pk comment('操作记录id') INT(11)"`
	ManagerId   int       `xorm:"not null comment('操作员id 对应员工id') INT(11)"`
	ManagerName string    `xorm:"not null default '' comment('操作人名称') VARCHAR(32)"`
	OperBiz     string    `xorm:"not null comment('操作业务 枚举') VARCHAR(64)"`
	OperParams  string    `xorm:"not null comment('操作内容参数') TEXT"`
	Created     time.Time `xorm:"not null default '1999-12-31 00:00:00' comment('操作时间') TIMESTAMP"`
	Updated     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status      int       `xorm:"not null default 1 comment('状态 1=正常') TINYINT(4)"`
}

func (t OperLog) TableName() string {
	return "oper_log"
}

func (t *OperLog) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

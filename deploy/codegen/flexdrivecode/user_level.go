package flexdrivecode

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type UserLevel struct {
	Id         int       `xorm:"not null pk autoincr comment('会员级别id') INT(10)"`
	Name       string    `xorm:"not null comment('会员级别名称') VARCHAR(32)"`
	QuotaSpace int64     `xorm:"not null default 0 comment('会员级别的用户空间 单位KB') BIGINT(20)"`
	Price      int       `xorm:"not null default 0 comment('会员级别的价格 单位分') INT(11)"`
	Created    time.Time `xorm:"not null default '1999-12-31 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status     int       `xorm:"not null default 1 comment('状态 1=启用 99=删除') TINYINT(4)"`
}

func (t UserLevel) TableName() string {
	return "user_level"
}

func (t *UserLevel) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type UserLevel struct {
	Id         int       `xorm:"not null pk autoincr comment('会员级别id') INT(10)"`
	Name       string    `xorm:"not null comment('会员级别名称') VARCHAR(32)"`
	QuotaSpace int64     `xorm:"not null default 0 comment('会员级别的用户空间 单位KB') BIGINT(20)"`
	Price      int       `xorm:"not null default 0 comment('会员级别的价格 单位分') INT(11)"`
	IsDefault  int       `xorm:"not null default 0 comment('是否是新用户的默认等级 1=是') TINYINT(4)"`
	Created    time.Time `xorm:"created not null default '1999-12-31 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated    time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status     int       `xorm:"not null default 1 comment('状态 1=启用 99=删除') TINYINT(4)"`
	PrimeCost  int       `xorm:"not null default 0 comment('原价 仅用于展示 单位分') INT(11)"`
	Desc       string    `xorm:"not null comment('介绍') VARCHAR(256)"`
}

func (t UserLevel) TableName() string {
	return "user_level"
}

func (t *UserLevel) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func GetUserLevelById(id int) (*UserLevel, error) {
	e := &UserLevel{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func GetUserLevelByName(name string) (*UserLevel, error) {
	e := &UserLevel{}
	exist, err := base.GetByCol("name", name, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func GetDefaultUserLevel() (*UserLevel, error) {
	e := &UserLevel{}
	exist, err := base.GetByCol("isdefault", 1, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

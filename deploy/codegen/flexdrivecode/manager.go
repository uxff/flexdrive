package flexdrivecode

import (
	"github.com/uxff/flexdrive/pkg/dao/base"
	"time"
)

type Manager struct {
	Id          int       `xorm:"not null pk autoincr comment('管理员id') INT(10)"`
	Name        string    `xorm:"not null default '' comment('管理员名称') VARCHAR(32)"`
	Phone       string    `xorm:"not null default '' comment('管理员手机号') VARCHAR(12)"`
	Email       string    `xorm:"not null default '' comment('管理员email') VARCHAR(32)"`
	Pwd         string    `xorm:"not null default '' comment('密码') VARCHAR(32)"`
	Created     time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status      int       `xorm:"not null default 1 comment('状态 1=正常 99=删除') TINYINT(4)"`
	RoleId      int       `xorm:"not null default 0 comment('角色id') INT(11)"`
	IsSuper     int       `xorm:"not null default 0 comment('是否是超管 1=超管') TINYINT(4)"`
	LastLoginAt time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('最后登录时间') TIMESTAMP"`
	LastLoginIp string    `xorm:"not null default '' comment('最后登录ip') VARCHAR(16)"`
}

func (t Manager) TableName() string {
	return "manager"
}

func (t *Manager) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

package dao

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/uxff/flexdrive/pkg/log"

	"github.com/uxff/flexdrive/pkg/dao/base"
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

// func (t *Manager) GetById(id int) error {
// 	_, err := base.GetByCol("id", id, t)
// 	return err
// }

func (t *Manager) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

// func (t *Manager) GetByEmail(email string) error {
// 	_, err := base.GetByCol("email", email, t)
// 	return err
// }

func (t *Manager) IsPwdValid(p string) bool {
	enc := md5.New()
	enc.Write([]byte(p))

	s := hex.EncodeToString(enc.Sum(nil))
	log.Debugf("IsPwdValid: p:%s p.md5:%s expected t.md5:%s t==s:%v", p, s, t.Pwd, s == t.Pwd)
	return s == t.Pwd
}
func (t *Manager) SetPwd(p string) {
	enc := md5.New()
	enc.Write([]byte(p))

	t.Pwd = hex.EncodeToString(enc.Sum(nil))
}

func (t *Manager) IsSuperRole() bool {
	return t.IsSuper == 1
}

func GetManagerById(id int) (*Manager, error) {
	m := &Manager{}
	_, err := base.GetByCol("id", id, m)
	return m, err
}

func GetManagerByEmail(email string) (*Manager, error) {
	m := &Manager{}
	_, err := base.GetByCol("email", email, m)
	return m, err
}

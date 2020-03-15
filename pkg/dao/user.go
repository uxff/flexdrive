package dao

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

type User struct {
	Id          int       `xorm:"not null pk autoincr comment('用户id') INT(10)"`
	Name        string    `xorm:"not null default '' comment('用户姓名') VARCHAR(32)"`
	Email       string    `xorm:"not null comment('邮箱') VARCHAR(32)"`
	Phone       string    `xorm:"not null default '' comment('手机号 ') VARCHAR(12)"`
	Pwd         string    `xorm:"not null default '' comment('密码') VARCHAR(32)"`
	LevelId     int       `xorm:"not null default 0 comment('级别id') INT(11)"`
	TotalCharge int       `xorm:"not null default 0 comment('累计充值 单位分') INT(11)"`
	QuotaSpace  int64     `xorm:"not null default 0 comment('当前拥有的空间 单位KB') BIGINT(20)"`
	UsedSpace   int64     `xorm:"not null default 0 comment('当前已用空间 单位KB') BIGINT(20)"`
	FileCount   int64     `xorm:"not null default 0 comment('文件数量') BIGINT(20)"`
	LastLoginAt time.Time `xorm:"default '0000-00-00 00:00:00' comment('最后登录时间') TIMESTAMP"`
	LastLoginIp string    `xorm:"not null default '' comment('最后登录ip') VARCHAR(16)"`
	Created     time.Time `xorm:"created not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated     time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status      int       `xorm:"not null default 1 comment('状态 1=正常启用 99=账户冻结 ') TINYINT(4)"`
}

func (t User) TableName() string {
	return "user"
}

func (t *User) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func (t *User) IsPwdValid(p string) bool {
	enc := md5.New()
	enc.Write([]byte(p))

	s := hex.EncodeToString(enc.Sum(nil))
	log.Debugf("IsPwdValid: p:%s p.md5:%s expected t.md5:%s t==s:%v", p, s, t.Pwd, s == t.Pwd)
	return s == t.Pwd
}

func (t *User) SetPwd(p string) {
	enc := md5.New()
	enc.Write([]byte(p))

	t.Pwd = hex.EncodeToString(enc.Sum(nil))
}

func (t *User) GetUserLevel() (ul *UserLevel) {
	if t.LevelId > 0 {
		ul, _ = GetUserLevelById(t.LevelId)
		return ul
	}
	ul, _ = GetDefaultUserLevel()
	if ul == nil {
		ul = &UserLevel{
			Id:   0,
			Name: "普通等级",
		}
	}
	return ul
}

func GetUserById(id int) (*User, error) {
	e := &User{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func GetUserByEmail(email string) (*User, error) {
	e := &User{}
	exist, err := base.GetByCol("email", email, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

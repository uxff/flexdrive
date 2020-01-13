package dao

import (
	"github.com/uxff/flexdrive/pkg/dao/base"
	"time"
)

type Role struct {
	Id      int       `xorm:"not null pk autoincr comment('角色id') INT(10)"`
	Name    string    `xorm:"not null default '' comment('角色名称') VARCHAR(32)"`
	Status  int       `xorm:"not null default 1 comment('状态 1=正常启用 99=删除') TINYINT(4)"`
	Created time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Permit  string    `xorm:"not null comment('授权内容 json') TEXT"`
}

func (t Role) TableName() string {
	return "role"
}

func (t *Role) GetById(id int) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

func (t *Role) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/app/admin/model/rbac"
	"github.com/uxff/flexdrive/pkg/dao/base"
)

type Role struct {
	Id      int             `xorm:"not null pk autoincr comment('角色id') INT(10)"`
	Name    string          `xorm:"not null default '' comment('角色名称') VARCHAR(32)"`
	Status  int             `xorm:"not null default 1 comment('状态 1=正常启用 99=删除') TINYINT(4)"`
	Created time.Time       `xorm:"created not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated time.Time       `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Permit  rbac.RoleAccess `xorm:"not null comment('授权内容 json') json"`
}

func (t Role) TableName() string {
	return "role"
}

func (t *Role) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func GetRoleById(id int) (*Role, error) {
	ent := &Role{}
	exist, err := base.GetByCol("id", id, ent)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return ent, nil
}
func GetRoleByName(name string) (*Role, error) {
	ent := &Role{}
	exist, err := base.GetByCol("name", name, ent)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return ent, err
}

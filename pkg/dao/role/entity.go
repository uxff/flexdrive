package role

import (
	"github.com/uxff/flexdrive/pkg/utils"
)

type Role struct {
	Rid      int32         `xorm:"pk autoincr COMMENT('角色id') 'rid'"`
	RName    string        `xorm:"varchar(40) COMMENT('角色名') 'rname'"`
	RCreated utils.JsonTime `xorm:"created not null default '0000-00-00 00:00:00' 'rcreated'"`
	RUpdated utils.JsonTime `xorm:"updated not null default '0000-00-00 00:00:00' 'rupdated'"`
	RDeleted utils.JsonTime `xorm:"not null default '0000-00-00 00:00:00' 'rdeleted'"`
	RStatus  int8          `xorm:"not null tinyint  default(1) COMMENT('状态') 'rstatus'"`
	RIsSuper int8          `xorm:"not null tinyint  default(0) COMMENT('是否是超管') 'rissuper'"`
	RAccess  string        `xorm:"text COMMENT('角色权限配置 json') 'raccess'"`
}

func (m Role) TableName() string {
	return "role"
}

func (m Role) IsSuper() bool {
	return m.RIsSuper > 0
	// return m.Rid < base.MaxSuperRoleId
}

func (entityPtr *Role) FillEmpty() {
	// pg 和 mysql 高版本要求时间格式不能为 0000-00-00
	if entityPtr.RDeleted.String() < "1000" {
		entityPtr.RDeleted.UnmarshalJSON([]byte("0001-01-01 00:00:01"))
	}
	if entityPtr.RUpdated.String() < "1000" {
		entityPtr.RUpdated.UnmarshalJSON([]byte("0001-01-01 00:00:01"))
	}
}

// xorm 支持的事件 校验和填充空数据
func (entityPtr *Role) BeforeInsert() {
	entityPtr.FillEmpty()
}

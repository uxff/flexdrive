package user

import (
	"github.com/uxff/flexdrive/pkg/utils"
)

type User struct {
	Id          int            `xorm:"pk autoincr COMMENT('用户id') "`
	Name        string         `xorm:"varchar(40) COMMENT('名称') "`
	Email       string         `xorm:"varchar(40) COMMENT('角色名') "`
	Phone       string         `xorm:"varchar(40) COMMENT('角色名') "`
	Pwd         string         `xorm:"varchar(40) COMMENT('角色名') "`
	LevelId     int            `xorm:"varchar(40) COMMENT('角色名') "`
	TotalCharge int            `xorm:"varchar(40) COMMENT('角色名') "`
	QuotaSpace  int            `xorm:"varchar(40) COMMENT('角色名') "`
	UsedSpace   int            `xorm:"varchar(40) COMMENT('角色名') "`
	Created     utils.JsonTime `xorm:"created not null default '0000-00-00 00:00:00'"`
	Updated     utils.JsonTime `xorm:"updated not null default '0000-00-00 00:00:00'"`
	Status      int8           `xorm:"not null tinyint  default(1) COMMENT('状态') "`
}

func (m User) TableName() string {
	return "role"
}

func (entityPtr *User) FillEmpty() {
	// pg 和 mysql 高版本要求时间格式不能为 0000-00-00
	if entityPtr.Updated.String() < "1000" {
		entityPtr.Updated.UnmarshalJSON([]byte("0001-01-01 00:00:01"))
	}
}

// xorm 支持的事件 校验和填充空数据
func (entityPtr *User) BeforeInsert() {
	entityPtr.FillEmpty()
}

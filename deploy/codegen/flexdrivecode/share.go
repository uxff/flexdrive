package flexdrivecode

import (
	"github.com/uxff/flexdrive/pkg/dao/base"
	"time"
)

type Share struct {
	Id       int       `xorm:"not null pk autoincr comment('文件id') INT(10)"`
	FileHash string    `xorm:"not null default '' comment('文件哈希') VARCHAR(32)"`
	UserId   int       `xorm:"not null default 0 comment('分享者用户id') INT(11)"`
	NodeId   int       `xorm:"not null default 0 comment('所在节点名') INT(11)"`
	FileName string    `xorm:"not null default '' comment('文件名') VARCHAR(32)"`
	Path     string    `xorm:"not null default '' comment('文件路径') VARCHAR(256)"`
	Created  time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status   int       `xorm:"not null default 1 comment('状态 1=正常 2=隐藏 99=已删除') TINYINT(4)"`
	Expired  time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('分享有效期') TIMESTAMP"`
}

func (t Share) TableName() string {
	return "share"
}

func (t *Share) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

package flexdrivecode

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type UserFile struct {
	Id       int       `xorm:"not null pk autoincr comment('文件id') INT(10)"`
	UserId   string    `xorm:"not null default '' comment('用户id') unique(IDX_userId_pathhash) VARCHAR(32)"`
	Path     string    `xorm:"not null default '' comment('文件路径') VARCHAR(256)"`
	FileName string    `xorm:"not null default '' comment('文件名') VARCHAR(256)"`
	PathHash string    `xorm:"not null comment('路径哈希，hash(path+fileName)，用户下唯一') unique(IDX_userId_pathhash) VARCHAR(32)"`
	FileHash string    `xorm:"not null default '' comment('文件哈希') VARCHAR(32)"`
	NodeId   int       `xorm:"not null default 0 comment('所在节点名 第一副本所在节点') INT(11)"`
	IsDir    int       `xorm:"not null default 0 comment('是否是目录') TINYINT(4)"`
	Created  time.Time `xorm:"not null default '1999-12-31 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status   int       `xorm:"not null default 1 comment('状态 1=正常 2=隐藏 99=下架') TINYINT(4)"`
	Size     int       `xorm:"not null default 0 comment('大小 单位Byte 目录则记录0') INT(11)"`
	Space    int       `xorm:"not null default 0 comment('占用空间单位 单位KB 目录则记录0') INT(11)"`
	Desc     string    `xorm:"not null comment('描述信息') TEXT"`
}

func (t UserFile) TableName() string {
	return "user_file"
}

func (t *UserFile) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

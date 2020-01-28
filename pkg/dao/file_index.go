package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type FileIndex struct {
	Id         int       `xorm:"not null pk autoincr comment('文件id') INT(10)"`
	FileName   string    `xorm:"not null default '' comment('文件名') VARCHAR(32)"`
	FileHash   string    `xorm:"not null default '' comment('文件内容哈希') VARCHAR(32)"`
	NodeId     int       `xorm:"not null default 0 comment('所在节点名 第一副本所在节点') INT(11)"`
	NodeId2    int       `xorm:"not null default 0 comment('所在节点名 第二副本所在节点') INT(11)"`
	NodeId3    int       `xorm:"not null default 0 comment('所在节点名 第三副本所在节点') INT(11)"`
	FilePath   string    `xorm:"not null default '' comment('文件文件') VARCHAR(256)"`
	Created    time.Time `xorm:"created not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated    time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status     int       `xorm:"not null default 1 comment('状态 0=未就绪 1=就绪 98=上传失败 99=删除') TINYINT(4)"`
	ReferCount int       `xorm:"not null default 0 comment('被引用数量') INT(11)"`
	Size       int       `xorm:"not null default 0 comment('大小 单位Byte') INT(11)"`
	Space      int       `xorm:"not null default 0 comment('占用空间单位 单位KB') INT(11)"`
	Desc       string    `xorm:"not null comment('描述信息') TEXT"`
}

func (t FileIndex) TableName() string {
	return "file_index"
}

func (t *FileIndex) GetById(id int) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

func (t *FileIndex) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

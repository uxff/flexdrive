package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/common"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/envinit"
	"github.com/uxff/flexdrive/pkg/log"
)

type FileIndex struct {
	Id         int       `xorm:"not null pk autoincr comment('文件索引id') INT(10)"`
	FileName   string    `xorm:"not null default '' comment('文件名 无用') VARCHAR(64)"`
	FileHash   string    `xorm:"not null default '' comment('文件内容哈希 唯一') VARCHAR(40)"` // 唯一
	NodeId     int       `xorm:"not null default 0 comment('所在节点名 第一副本所在节点') INT(11)"`
	NodeId2    int       `xorm:"not null default 0 comment('所在节点名 第二副本所在节点') INT(11)"`
	NodeId3    int       `xorm:"not null default 0 comment('所在节点名 第三副本所在节点') INT(11)"`
	InnerPath  string    `xorm:"not null default '' comment('文件本地路径 在节点上的实际路径') VARCHAR(256)"`
	OuterPath  string    `xorm:"not null default '' comment('文件外部访问路径') VARCHAR(256)"`
	Created    time.Time `xorm:"created not null default '1999-12-31 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated    time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status     int       `xorm:"not null default 1 comment('状态 0=未就绪 1=就绪 98=上传失败 99=删除') TINYINT(4)"`
	ReferCount int       `xorm:"not null default 0 comment('被引用数量') INT(11)"`
	Size       int64     `xorm:"not null default 0 comment('大小 单位Byte') BIGINT(20)"`
	Space      int64     `xorm:"not null default 0 comment('占用空间单位 单位KB') BIGINT(20)"`
	Desc       string    `xorm:"not null comment('描述信息 无用') TEXT"`
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

func GetFileIndexById(id int) (*FileIndex, error) {
	e := &FileIndex{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func GetFileIndexByFileHash(fileHash string) (*FileIndex, error) {
	e := &FileIndex{}
	exist, err := base.GetByCol("fileHash", fileHash, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

// 统计文件数
func (t *FileIndex) CountRefers() int64 {
	dbname := common.DBMysqlDrive // t.DbNamespace()
	// total, err := engine.Where("id >?", 1).SumInt(ss, "money")
	total, err := envinit.Dbs[dbname].Where("fileIndexId=? and status=?", t.Id, 1).Count(&UserFile{}, "id")

	if err != nil {
		log.Errorf("count fileIndex(%d).files failed:%v", t.Id, err)
		return 0
	}
	t.ReferCount = int(total)
	err = t.UpdateById([]string{"referCount"})
	if err != nil {
		log.Errorf("count fileIndex(%d).files failed:%v", t.Id, err)
		return 0
	}
	return total
}

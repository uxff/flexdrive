package dao

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

type UserFile struct {
	Id          int       `xorm:"not null pk autoincr comment('用户文件索引id') INT(10)"`
	FileIndexId int       `xorm:"not null default '0' comment('文件索引id') INT(11)"`
	UserId      int       `xorm:"not null default '0' comment('用户id') index(IDX_userId_pathhash) int(11)"`
	FilePath    string    `xorm:"not null default '' comment('文件路径 用户展示路径 ') VARCHAR(256)"`
	FileName    string    `xorm:"not null default '' comment('文件名 如果是目录则空 用户可改名') VARCHAR(256)"`
	PathHash    string    `xorm:"not null comment('路径哈希，hash(filePath)') index(IDX_userId_pathhash) VARCHAR(40)"`
	FileHash    string    `xorm:"not null default '' comment('文件哈希 如果是目录则空') VARCHAR(32)"`
	NodeId      int       `xorm:"not null default 0 comment('所在节点名 第一副本所在节点') INT(11)"`
	IsDir       int       `xorm:"not null default 0 comment('是否是目录') TINYINT(4)"`
	Created     time.Time `xorm:"created not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated     time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status      int       `xorm:"not null default 1 comment('状态 1=正常 2=隐藏 99=下架') TINYINT(4)"`
	Size        int       `xorm:"not null default 0 comment('大小 单位Byte 目录则记录0') INT(11)"`
	Space       int       `xorm:"not null default 0 comment('占用空间单位 单位KB 目录则记录0') INT(11)"`
	Desc        string    `xorm:"not null comment('描述信息') TEXT"`

	// after select
	User      *User      `xorm:"-"`
	FileIndex *FileIndex `xorm:"-"`
}

func (t UserFile) TableName() string {
	return "user_file"
}

func (t *UserFile) GetById(id int) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

func (t *UserFile) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func (t *UserFile) MakePathHash() string {
	encSha1 := sha1.New()
	encSha1.Write([]byte(t.FilePath))
	t.PathHash = hex.EncodeToString(encSha1.Sum([]byte("")))
	return t.PathHash
}

func GetUserFileById(id int) (*UserFile, error) {
	e := &UserFile{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func (t *UserFile) AfterSelect() {
	var err error
	t.User, err = GetUserById(t.UserId)
	if err != nil {
		log.Warnf("load userFile.User error:%v", err)
	}

	t.FileIndex, err = GetFileIndexById(t.FileIndexId)
	if err != nil {
		log.Warnf("load userFile.FileIndex error:%v", err)
	}

	log.Debugf("load userFile.User ok")
}

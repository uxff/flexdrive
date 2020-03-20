package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/log"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

const (
	OfflineTaskStatusExecuting = 1
	//OfflineTaskStatusDone      = 2
	OfflineTaskStatusFail    = 3
	OfflineTaskStatusSaved   = 4
	OfflineTaskStatusRemoved = 99
)

type OfflineTask struct {
	Id               int       `xorm:"not null pk autoincr comment('主键id') INT(10)"`
	UserId           int       `xorm:"not null default 0 comment('用户id') INT(11)"`
	Dataurl          string    `xorm:"not null default '' comment('资源地址') VARCHAR(256)"`
	FileName         string    `xorm:"not null default '' comment('文件名') VARCHAR(64)"`
	Created          time.Time `xorm:"created not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated          time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status           int       `xorm:"not null default 1 comment('状态 1=下载中 2=下载完成 3=下载失败 4=已保存') TINYINT(4)"`
	ParentUserFileId int       `xorm:"not null default 0 comment('文件索引id') INT(11)"`
	UserFileId       int       `xorm:"not null default 0 comment('文件索引id') INT(11)"`
	FileHash         string    `xorm:"not null default '' comment('文件哈希') VARCHAR(40)"`
	Size             int64     `xorm:"not null default 0 comment('大小') BIGINT(20)"`
	Remark           string    `xorm:"not null default '' comment('备注 比如失败原因') VARCHAR(256)"`

	// after select
	User     *User     `xorm:"-"`
	UserFile *UserFile `xorm:"-"`
}

func (t OfflineTask) TableName() string {
	return "offline_task"
}

func (t *OfflineTask) GetById(id int) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

func (t *OfflineTask) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func GetOfflineTaskById(id int) (*OfflineTask, error) {
	e := &OfflineTask{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

//
func (t *OfflineTask) AfterSelect() {
	var err error
	t.User, err = GetUserById(t.UserId)
	if err != nil {
		log.Warnf("load OfflineTask(%d).User(%d) error:%v", t.Id, t.UserId, err)
	}
	t.UserFile, err = GetUserFileById(t.UserFileId)
	if err != nil {
		log.Warnf("load OfflineTask(%d).UserFile(%d) error:%v", t.Id, t.UserFileId, err)
	}

	log.Debugf("load OfflineTask %d .User .UserFile ok", t.Id)
}

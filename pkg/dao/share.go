package dao

import (
	"fmt"
	"time"

	"github.com/uxff/flexdrive/pkg/utils/filehash"

	"github.com/uxff/flexdrive/pkg/log"

	"github.com/uxff/flexdrive/pkg/dao/base"
)

type Share struct {
	Id         int       `xorm:"not null pk autoincr comment('文件id') INT(10)"`
	ShareHash  string    `xorm:"not null default '' comment('share哈希 用于访问 带索引') VARCHAR(40)"`
	FileHash   string    `xorm:"not null default '' comment('文件哈希') VARCHAR(40)"`
	UserId     int       `xorm:"not null default 0 comment('分享者用户id') INT(11)"`
	UserFileId int       `xorm:"not null default 0 comment('分享者文件索引id') INT(11)"`
	NodeId     int       `xorm:"not null default 0 comment('所在节点id') INT(11)"`
	FileName   string    `xorm:"not null default '' comment('文件名') VARCHAR(64)"`
	Created    time.Time `xorm:"created not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated    time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status     int       `xorm:"not null default 1 comment('状态 1=正常 2=隐藏 99=已删除') TINYINT(4)"`
	Expired    time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('分享有效期') TIMESTAMP"`

	// after select
	User     *User     `xorm:"-"`
	UserFile *UserFile `xorm:"-"`
}

func (t Share) TableName() string {
	return "share"
}

func (t *Share) GetById(id int) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

func (t *Share) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

//  保证UserFileId已经赋值
func (t *Share) MakeShareHash() string {
	raw := fmt.Sprintf("share-%d", t.UserFileId)
	t.ShareHash, _ = filehash.CalcSha1(raw)
	return t.ShareHash
}

func GetShareById(id int) (*Share, error) {
	e := &Share{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func (t *Share) AfterSelect() {
	var err error
	t.User, err = GetUserById(t.UserId)
	if err != nil {
		log.Warnf("load share.User error:%v", err)
	}
	t.UserFile, err = GetUserFileById(t.UserFileId)
	if err != nil {
		log.Warnf("load share.UserFile error:%v", err)
	}

	log.Debugf("load share.User, share.UserFile ok")
}

func GetShareByUserFile(userFileId int) (*Share, error) {
	e := &Share{}
	exist, err := base.GetByCol("userFileId", userFileId, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, nil
}

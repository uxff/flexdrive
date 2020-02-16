package dao

import (
	"crypto/sha1"
	"encoding/hex"
	"path"
	"strings"
	"time"

	"github.com/uxff/flexdrive/pkg/common"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/envinit"
	"github.com/uxff/flexdrive/pkg/log"
)

type UserFile struct {
	Id          int       `xorm:"not null pk autoincr comment('用户文件索引id') INT(10)"`
	FileIndexId int       `xorm:"not null default '0' comment('文件索引id') INT(11)"`
	UserId      int       `xorm:"not null default '0' comment('用户id') index(IDX_userId_pathhash) int(11)"`
	FilePath    string    `xorm:"not null default '' comment('文件父路径 用户展示路径 /结尾 ') VARCHAR(256)"`
	FileName    string    `xorm:"not null default '' comment('文件名 如果是目录则空 用户可改名') VARCHAR(256)"`
	PathHash    string    `xorm:"not null comment('路径哈希，hash(filePath)') index(IDX_userId_pathhash) VARCHAR(40)"`
	FileHash    string    `xorm:"not null default '' comment('文件内容哈希 如果是目录则空') VARCHAR(40)"`
	NodeId      int       `xorm:"not null default 0 comment('所在节点名 第一副本所在节点') INT(11)"`
	IsDir       int       `xorm:"not null default 0 comment('是否是目录') TINYINT(4)"`
	Created     time.Time `xorm:"created not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated     time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Status      int       `xorm:"not null default 1 comment('状态 1=正常 2=隐藏 99=下架') TINYINT(4)"`
	Size        int64     `xorm:"not null default 0 comment('大小 单位Byte 目录则记录0') BIGINT(20)"`
	Space       int64     `xorm:"not null default 0 comment('占用空间单位 单位KB 目录则记录0') BIGINT(20)"`
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

// FilePath指的是文件父目录
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

// 获取目录一级的UserFile对象 查看一个对象是否存在
// 如果是目录 则查看 t.filePath=path.Dir(filePath) and t.FileName=path.Base(filePath)
// 如果filePath=/abc/ 则查看/abc这个目录是否存在
// 如果filePath=/abc/efg 则查看/abc/下efg这个目录或者文件是否存在
func GetUserFileByPath(userId int, filePath string) (*UserFile, error) {
	if filePath == "/" {
		return &UserFile{
			FilePath: "/",
		}, nil
	}

	// 在本系统中filePath要以/结尾
	parentPath := path.Dir(strings.TrimRight(filePath, "/"))

	if len(parentPath) > 1 {
		parentPath += "/"
	}
	baseName := path.Base(strings.TrimRight(filePath, "/"))
	log.Debugf("filePath:%s parentPath:%s baseName:%s", filePath, parentPath, baseName)

	e := &UserFile{
		FilePath: parentPath,
	}

	list := make([]*UserFile, 0)
	err := base.ListByCondition(e, map[string]interface{}{
		"userId=?":   userId,
		"pathHash=?": e.MakePathHash(),
		"filePath=?": parentPath,
		"fileName=?": baseName,
		"status=?":   base.StatusNormal,
	}, 1, 1000, "id", &list)

	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	item := list[0]
	return item, err
}

// 创建目录一级的UserFile对象 filePath需要自己补充
func MakeUserFilePath(userId int, filePath, name string) (*UserFile, error) {
	e := &UserFile{
		FilePath: filePath,
		FileName: name,
		IsDir:    1,
		Status:   1,
		UserId:   userId,
	}

	e.MakePathHash()

	_, err := base.Insert(e)
	if err != nil {
		return nil, err
	}
	return e, nil
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

// 统计用户的使用空间
func (t *UserFile) SumSpace() int64 {
	dbname := common.DBMysqlDrive // t.DbNamespace()
	// if entityOfDb, ok := t.(base.DbNamespace); ok {
	// 	dbname = entityOfDb.DbNamespace()
	// }

	// total, err := engine.Where("id >?", 1).SumInt(ss, "money")
	total, err := envinit.Dbs[dbname].Where("userId=? and status=?", t.UserId, t.Status).SumInt(&UserFile{}, "space")

	if err != nil {
		log.Errorf("sum user(%d).usedSpace failed:%v", t.UserId, err)
		return 0
	}
	return total
}

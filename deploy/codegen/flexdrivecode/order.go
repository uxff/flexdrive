package flexdrivecode

import (
	"github.com/uxff/flexdrive/pkg/dao/base"
	"time"
)

type Order struct {
	Id            int       `xorm:"not null pk autoincr comment('订单id') INT(10)"`
	UserId        int       `xorm:"not null comment('用户id') INT(11)"`
	OriginLevelId int       `xorm:"not null comment('用户原等级id') INT(11)"`
	AwardLevelId  int       `xorm:"not null default 0 comment('用户购买的等级id') INT(11)"`
	AwardSpace    int64     `xorm:"not null default 0 comment('本次购买的容量空间 单位KB') BIGINT(11)"`
	Phone         string    `xorm:"not null default '' comment('用户手机号') VARCHAR(12)"`
	LevelName     string    `xorm:"not null default '' comment('等级名') VARCHAR(12)"`
	TotalAmount   int       `xorm:"not null default 0 comment('订单价格 单位分') INT(11)"`
	PayAmount     int       `xorm:"not null default 0 comment('实付款金额 单位分') INT(11)"`
	Remark        string    `xorm:"not null comment('订单备注') TEXT"`
	Created       time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	Status        int       `xorm:"not null default 1 comment('状态 1=待付款 2=未付款关闭 3=已付款 4=已退款') TINYINT(4)"`
}

func (t Order) TableName() string {
	return "order"
}

func (t *Order) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

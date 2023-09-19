package dao

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

const (
	OrderStatusPaying    = 1
	OrderStatusClosed    = 2
	OrderStatusPaid      = 3
	OrderStatusRefunding = 4
	OrderStatusRefended  = 5
)

var OrderStatusMap = map[int]string{
	OrderStatusPaying:    "待付款",
	OrderStatusClosed:    "未付款关闭",
	OrderStatusPaid:      "支付完成",
	OrderStatusRefunding: "退款中",
	OrderStatusRefended:  "退款完成",
}

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
	OutOrderNo    string    `xorm:"not null comment('第三方支付通道订单号') VARCHAR(40)"`
	Remark        string    `xorm:"not null comment('订单备注') TEXT"`
	Created       time.Time `xorm:"created not null default '1999-12-31 00:00:00' comment('创建时间') TIMESTAMP"`
	Updated       time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	Status        int       `xorm:"not null default 1 comment('状态 1=待付款 2=未付款关闭 3=已付款 4=退款中 5=已退款') TINYINT(4)"`

	// after select
	User *User `xorm:"-"`
}

func (t Order) TableName() string {
	return "order"
}

func (t *Order) GetById(id int) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

func (t *Order) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func GetOrderById(id int) (*Order, error) {
	e := &Order{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

func (t *Order) AfterSelect() {
	var err error
	t.User, err = GetUserById(t.UserId)
	if err != nil {
		log.Warnf("load order.User error:%v", err)
	}

	log.Debugf("load order.User ok")
}

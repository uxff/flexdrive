package user

import (
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

func GetById(id int32) (*User, error) {
	e := User{}

	//has, err := base.GetByColWithCache("id", id, &e)
	has, err := base.GetByCol("id", id, &e)
	// has, err := db.Dbs[common.DBNameNamespace].Engine.Where("merAppId=?", merAppId).Get(&mc)
	if err != nil {
		log.Errorf("query(id=%d) failed:%v", id, err)
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return &e, nil
}

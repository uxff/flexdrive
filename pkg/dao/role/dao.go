package role

import (
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

func QueryByRid(rid int32) (*Role, error) {
	m := Role{}

	has, err := base.GetByColWithCache("rid", rid, &m)
	// has, err := db.Dbs[common.DBNameNamespace].Engine.Where("merAppId=?", merAppId).Get(&mc)
	if err != nil {
		log.Errorf("query(rid=%d) failed:%v", rid, err)
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return &m, nil
}

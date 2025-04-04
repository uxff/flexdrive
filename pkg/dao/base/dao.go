package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/uxff/flexdrive/pkg/common"
	"github.com/uxff/flexdrive/pkg/envinit"
	"github.com/uxff/flexdrive/pkg/log"
)

// 是否使用缓存
const (
	DontUseCache  = 0
	UseRedisCache = 1
)

type TableName interface {
	TableName() string
}

type DbNamespace interface {
	DbNamespace() string
}

type AfterSelect interface {
	AfterSelect()
}

// 缓存过期时间 30min
var cacheExpireSec int64 = 60 * 30

func Insert(entityPtr interface{}) (int64, error) {

	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityPtr.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	tsStart := time.Now()

	session := envinit.Dbs[dbname].NewSession()
	n, err := session.Insert(entityPtr)

	// sql, params := session.LastSQL() // not work
	log.Debugf("params:%+v timeused:%dms affected:%v", entityPtr, time.Since(tsStart).Nanoseconds()/1000/1000, n)

	if err != nil {
		log.Errorf("insert error:%v", err)
	}

	return n, err
}

/*
*

	    查找 colName=colVal 的记录 写入到entityPtr中
		colName 一般为表的唯一索引 比如mchconfigs.MchConfig.GetKeyName返回MchId
*/
func GetByCol(colName string, colVal interface{}, entityPtr interface{}) (found bool, err error) {
	// auto load cache
	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityPtr.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	tsStart := time.Now()

	session := envinit.Dbs[dbname].Where(colName+" = ?", colVal)
	found, err = session.Get(entityPtr)

	// sql, params := session.LastSQL() // not work
	log.Debugf("params:%+v timeused:%dms found:%v", colVal, time.Since(tsStart).Nanoseconds()/1000/1000, found)

	if v, ok := entityPtr.(AfterSelect); ok {
		v.AfterSelect()
	}
	return
}

/*
*

	    查找 where colName=colVal 的记录 写入到entityPtr中
		colName 一般为表的唯一索引 比如mchconfigs.MchConfig.GetKeyName返回MchId
		要保证传入的entityPtr不为空指针
*/
func GetByColWithCache(colName string, colVal interface{}, entityPtr interface{}) (found bool, err error) {
	// auto load cache
	v, ok := entityPtr.(TableName)
	if !ok {
		return false, errors.New("invalid entityPtr, expect bean")
	}

	cacheKey := v.TableName() + "-" + colName + "-" + fmt.Sprintf("%v", colVal)
	if CacheGet(cacheKey, entityPtr) == nil {
		log.Debugf("got from cache:" + cacheKey)
		return true, nil
	}

	found, err = GetByCol(colName, colVal, entityPtr)
	if found {
		CacheSet(cacheKey, entityPtr)
	}

	return
}

// 保证entityPtr不为空指针，本函数会对其进行写
func CacheGet(cacheKey string, entityPtr interface{}) error {
	redisConn, _ := envinit.GetRdsConn(common.RedisDrive)
	if redisConn == nil {
		return errors.New("get redis connection failed")
	}
	defer redisConn.Close()

	redisValue, err := redisConn.GetString(cacheKey)
	if err != nil {
		return err
	}

	if redisValue == "" {
		return errors.New("empty value in cache:" + cacheKey)
	}

	err = json.Unmarshal([]byte(redisValue), entityPtr)
	return err
}

// 确保entityPtr对应的对象中有值，本函数只读其值
func CacheSet(cacheKey string, entityPtr interface{}) error {
	redisConn, _ := envinit.GetRdsConn(common.RedisDrive)
	if redisConn == nil {
		return errors.New("get redis connection failed")
	}
	defer redisConn.Close()

	redisVal, err := json.Marshal(entityPtr)
	if err != nil {
		return err
	}

	_, err = redisConn.SetEx(cacheKey, cacheExpireSec, string(redisVal))
	return err
}

// entityEmptyPtr可空
func CacheDel(cacheKey string) error {
	redisConn, _ := envinit.GetRdsConn(common.RedisDrive)
	if redisConn == nil {
		return errors.New("get redis connection failed")
	}
	defer redisConn.Close()

	_, err := redisConn.Del(cacheKey)
	log.Debugf("cache deleted: key=%s err=%v", cacheKey, err)
	return err
}

/*
*

	更新： 查找 where colName=colVal 的记录 更新按照entityPtr中取cols指定的字段更新到数据库
*/
func UpdateByCol(colName string, colVal interface{}, entityPtrWithValue interface{}, cols []string) (n int64, err error) {

	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityPtrWithValue.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	tsStart := time.Now()

	session := envinit.Dbs[dbname].Cols(cols...).Where(colName+" = ?", colVal)
	n, err = session.Update(entityPtrWithValue)

	// sql, params := session.LastSQL() // not work
	log.Debugf("params:%+v timeused:%dms affected:%d", colVal, time.Since(tsStart).Nanoseconds()/1000/1000, n)

	return
}

// 更新 同UpdateByCol 额外操作了删除缓存
func UpdateByColWithCache(colName string, colVal interface{}, entityPtrWithValue TableName, cols []string) (n int64, err error) {
	// auto update cache
	n, err = UpdateByCol(colName, colVal, entityPtrWithValue, cols)

	// 删除对应的缓存
	go CacheDelByEntity(colName, colVal, entityPtrWithValue)
	return
}

/*
*

	软删除 主键字段值=key 的记录 字段值从entityPtr中取
*/
func DeleteByCol(colName string, colVal interface{}, statusColName string, entityEmptyPtr TableName) (err error) {
	// auto update cache
	// cols := []string{"status"}

	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityEmptyPtr.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	tsStart := time.Now()

	session := envinit.Dbs[dbname].Table(entityEmptyPtr).Cols(statusColName).Where(colName+" = ?", colVal)
	n, err := session.Update(map[string]interface{}{
		statusColName: StatusDeleted, // 此处转小写，否则pg不兼容
	})

	// sql, params := session.LastSQL() // not work
	log.Debugf("params:%+v timeused:%dms affected:%d", colVal, time.Since(tsStart).Nanoseconds()/1000/1000, n)

	// 删除对应的缓存
	go CacheDelByEntity(colName, colVal, entityEmptyPtr)
	return err
}

// 只过期缓存 不操作数据库
func CacheDelByEntity(colName string, colVal interface{}, entityPtr TableName) error {
	if entityPtr == nil {
		return nil
	}

	cacheKey := entityPtr.TableName() + "-" + colName + "-" + fmt.Sprintf("%v", colVal)
	err := CacheDel(cacheKey)
	if err != nil {
		log.Errorf("cache del %s error:%v", cacheKey, err)
		return err
	}

	return nil
}

/*
*
更新： 查找 where conditions 的记录 更新按照entityPtrWithValue中取cols指定的字段更新到数据库
conditions = ["a = ?"=>1,"b like '%?%'"=>"bb"]
// 允许conditions里key的value为空
*/
func UpdateByCondition(entityPtrWithValue interface{}, conditions map[string]interface{}, cols []string) (n int64, err error) {

	whereVal := make([]interface{}, 0)
	whereStr := "1 " // postgre 要求true开头; mysql 可以true或1; sqlite要求不能是true
	for ck, cv := range conditions {
		whereStr += " and " + ck
		if cv == nil {
			continue
		}
		whereVal = append(whereVal, cv)
	}

	tsStart := time.Now()
	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityPtrWithValue.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	session := envinit.Dbs[dbname].Cols(cols...).Where(whereStr, whereVal...)
	n, err = session.Update(entityPtrWithValue)

	// sql, params := session.LastSQL() // not work
	log.Debugf("params:%+v timeused:%dms affected:%d", whereVal, time.Since(tsStart).Nanoseconds()/1000/1000, n)

	return
}

/*
*
conditions = ["a = ?"=>1,"b like '%?%'"=>"bb"]
// 如果conditions map的key里不包含“?”，则默认认为value是slice,即使用where in
// 允许conditions里key的value为空
*/
func ListAndCountByCondition(entityPtr interface{}, conditions map[string]interface{}, pageNo int, pageSize int, orderBy string, listSlicePtr interface{}) (total int64, err error) {
	if pageNo <= 0 {
		pageNo = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	whereVal := make([]interface{}, 0)
	whereStr := "1 " // postgre 要求true开头; mysql 可以true或1; sqlite要求不能是true
	whererInMap := make(map[string]interface{})
	for ck, cv := range conditions {
		if cv == nil {
			whereStr += " and " + ck
			continue
		}
		if !strings.Contains(ck, "?") {
			whererInMap[ck] = cv
		} else {
			whereStr += " and " + ck
			whereVal = append(whereVal, cv)
		}
	}

	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityPtr.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	session := envinit.Dbs[dbname].Where(whereStr, whereVal...)
	for k, v := range whererInMap {
		session.In(k, v)
	}

	var nCount int64
	nCount, err = session.Count(entityPtr)
	if err != nil {
		return
	}

	//log.Debugf("nCount of listByCondition =%d", nCount)

	total = nCount

	tsStart := time.Now()

	// 此处必须重新建一个session,不能继续用上次的session
	session = envinit.Dbs[dbname].Where(whereStr, whereVal...)
	for k, v := range whererInMap {
		session.In(k, v)
	}
	err = session.OrderBy(orderBy).Limit(pageSize, (pageNo-1)*pageSize).Find(listSlicePtr)

	listOfPage := reflect.ValueOf(listSlicePtr).Elem()
	// sql, params := session.LastSQL() // not work
	log.Debugf("params:%+v timeused:%dms count:%d", whereVal, time.Since(tsStart).Nanoseconds()/1000/1000, reflect.ValueOf(listSlicePtr).Elem().Len())

	if _, ok := entityPtr.(AfterSelect); ok {

		//listSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(entityPtr))).Elem().Addr().Interface()

		// 处理页中的每行
		for i := 0; i < listOfPage.Len(); i++ {
			rowEntity := listOfPage.Index(i).Interface()
			if vv, vvok := rowEntity.(AfterSelect); vvok {
				vv.AfterSelect()
			}
		}

	}
	return
}

/*
*
conditions = ["a = ?"=>1,"b like '%?%'"=>"bb"]
// 允许conditions里key的value为空
// listSlicePtr for output slice
*/
func ListByCondition(entityPtr interface{}, conditions map[string]interface{}, pageNo int, pageSize int, orderBy string, listSlicePtr interface{}) (err error) {

	if pageNo <= 0 {
		pageNo = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	whereVal := make([]interface{}, 0)
	whereStr := "1 " // postgre 要求true开头; mysql 可以true或1; sqlite要求不能是true
	for ck, cv := range conditions {
		whereStr += " and " + ck
		if cv == nil {
			continue
		}
		whereVal = append(whereVal, cv)
	}

	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityPtr.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	tsStart := time.Now()

	session := envinit.Dbs[dbname].Where(whereStr, whereVal...).Limit(pageSize, (pageNo-1)*pageSize)
	if orderBy != "" {
		session = session.OrderBy(orderBy)
	}

	err = session.Find(listSlicePtr)

	// sql, params := session.LastSQL() // not work
	listOfPage := reflect.ValueOf(listSlicePtr).Elem()
	log.Debugf("params:%+v timeused:%dms count:%d", whereVal, time.Since(tsStart).Nanoseconds()/1000/1000, listOfPage.Len())

	if _, ok := entityPtr.(AfterSelect); ok {

		//listSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(entityPtr))).Elem().Addr().Interface()

		// 处理页中的每行
		for i := 0; i < listOfPage.Len(); i++ {
			rowEntity := listOfPage.Index(i).Interface()
			if vv, vvok := rowEntity.(AfterSelect); vvok {
				vv.AfterSelect()
			}
		}

	}

	return
}

// work well
// 如果 emptyEntity 是 *agents.Agent 则 eachRowHandler 中 emptyEntity 的类型是 *agents.Agent
// 如果 emptyEntity 是  agents.Agent 则 eachRowHandler 中 emptyEntity 的类型是  agents.Agent
func RangeTableByConditions(tableEntity interface{}, conditions map[string]interface{}, startPageNo, pageSize int, orderBy string, eachRowHandler func(rowNo int, rowEntity interface{})) (nCount int, err error) {
	if startPageNo <= 0 {
		startPageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 1000
	}

	if eachRowHandler == nil {
		return 0, errors.New("when RangeTableByConditions eachRowHandler cannot be nil")
	}

	tableName := getTableName(tableEntity)

	whereVal := make([]interface{}, 0)
	whereStr := "1 "
	for ck, cv := range conditions {
		whereStr += " and " + ck
		if cv == nil {
			continue
		}
		whereVal = append(whereVal, cv)
	}

	dbname := common.DBMysqlDrive
	if entityOfDb, ok := tableEntity.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	// 每页读取
	for {
		tsStart := time.Now()

		session := envinit.Dbs[dbname].Where(whereStr, whereVal...).Limit(pageSize, (startPageNo-1)*pageSize)
		if orderBy != "" {
			session = session.OrderBy(orderBy)
		}

		//listSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(entityPtr)), pageSize, pageSize)
		listSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(tableEntity))).Elem().Addr().Interface()

		err = session.Find(listSlice)

		// sql, params := session.LastSQL() // not work
		log.Debugf("params:%+v timeused:%dms count:%d", whereVal, time.Since(tsStart).Nanoseconds()/1000/1000, reflect.ValueOf(listSlice).Elem().Len())

		if err != nil {
			log.Errorf("RangeTableByConditions error:%v", err)
			return nCount, err
		}

		listOfPage := reflect.ValueOf(listSlice).Elem()
		log.Debugf("get page:%d count:%d/%d", startPageNo, listOfPage.Len(), pageSize)

		// 处理页中的每行
		for i := 0; i < listOfPage.Len(); i++ {
			rowEntity := listOfPage.Index(i).Interface()
			eachRowHandler(startPageNo*pageSize+i, rowEntity)
			//log.Debugf("got a item from range all::%+v", entityPtr)
			nCount++
		}
		//log.Debugf("pageNo:%d listSlice.Len:%d", startPageNo, reflect.ValueOf(listSlice).Elem().Len())

		if reflect.ValueOf(listSlice).Elem().Len() < pageSize {
			log.Debugf("RangeTableByConditions lastPage:%d allRow:%d tableEntity:%v", startPageNo, startPageNo*pageSize+reflect.ValueOf(listSlice).Elem().Len(), tableName)
			break
		}

		startPageNo++
	}

	log.Debugf("range table:%s done pageNo:%d allCount:%d", tableName, pageSize, nCount)

	return nCount, nil
}

/*
*
conditions = ["a = ?"=>1,"b like '%?%'"=>"bb"]
// 允许conditions里key的value为空
*/
func CountByCondition(entityPtr interface{}, conditions map[string]interface{}) (total int64, err error) {
	whereVal := make([]interface{}, 0)
	whereStr := "1 "
	for ck, cv := range conditions {
		whereStr += " and " + ck
		if cv == nil {
			continue
		}
		whereVal = append(whereVal, cv)
	}

	dbname := common.DBMysqlDrive
	if entityOfDb, ok := entityPtr.(DbNamespace); ok {
		dbname = entityOfDb.DbNamespace()
	}

	tsStart := time.Now()

	session := envinit.Dbs[dbname].Where(whereStr, whereVal...)
	total, err = session.Count(entityPtr)

	// sql, params := session.LastSQL() // not work
	log.Debugf("params:%+v timeused:%dms count:%d", whereVal, time.Since(tsStart).Nanoseconds()/1000/1000, total)

	// log.Debugf("nCount of listByCondition =%d", nCount)
	if err != nil {
		log.Errorf("countByCondition error:%v", err)
		return
	}
	return
}

func getTableName(tableEntity interface{}) (tableName string) {
	if t, ok := tableEntity.(TableName); ok {
		tableName = t.TableName()
	}

	if tableName != "" {
		return
	}

	var tableEntityPtr interface{} = &tableEntity
	if t, ok := tableEntityPtr.(TableName); ok {
		tableName = t.TableName()
	}

	return
}

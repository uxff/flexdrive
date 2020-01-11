package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uxff/flexdrive/pkg/common"
	"github.com/uxff/flexdrive/pkg/envinit"
	"github.com/uxff/flexdrive/pkg/log"
	"reflect"
	"strings"
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

// 缓存过期时间 30min
var cfgExpireSec int64 = 60 * 30

func Insert(entityPtr interface{}) (int64, error) {
	n, err := envinit.Dbs[common.DBNamespace].Insert(entityPtr)
	if err != nil {
		log.Errorf("insert error:%v", err)
	}

	return n, err
}

/**
    查找 colName=colVal 的记录 写入到entityPtr中
	colName 一般为表的唯一索引 比如mchconfigs.MchConfig.GetKeyName返回MchId
*/
func GetByCol(colName string, colVal interface{}, entityPtr interface{}) (found bool, err error) {
	// auto load cache
	//if common.DBNamespace == common.DBPostgre {
	//	colName = fmt.Sprintf(`"%s"`, colName)
	//}
	found, err = envinit.Dbs[common.DBNamespace].Where(colName+" = ?", colVal).Get(entityPtr)
	return
}

/**
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
	redisConn, _ := envinit.GetRdsConn(common.RedisNamespace)
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
	redisConn, _ := envinit.GetRdsConn(common.RedisNamespace)
	if redisConn == nil {
		return errors.New("get redis connection failed")
	}
	defer redisConn.Close()

	redisVal, err := json.Marshal(entityPtr)
	if err != nil {
		return err
	}

	_, err = redisConn.SetEx(cacheKey, cfgExpireSec, string(redisVal))
	return err
}

// entityEmptyPtr可空
func CacheDel(cacheKey string) error {
	redisConn, _ := envinit.GetRdsConn(common.RedisNamespace)
	if redisConn == nil {
		return errors.New("get redis connection failed")
	}
	defer redisConn.Close()

	_, err := redisConn.Del(cacheKey)
	log.Debugf("cache deleted: key=%s err=%v", cacheKey, err)
	return err
}

/**
  更新： 查找 where colName=colVal 的记录 更新按照entityPtr中取cols指定的字段更新到数据库
*/
func UpdateByCol(colName string, colVal interface{}, entityPtrWithValue interface{}, cols []string) (n int64, err error) {

	//if common.DBNamespace == common.DBPostgre {
	//	colName = fmt.Sprintf(`"%s"`, colName)
	//}
	n, err = envinit.Dbs[common.DBNamespace].Cols(cols...).Where(colName+" = ?", colVal).Update(entityPtrWithValue)
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

/**
  软删除 主键字段值=key 的记录 字段值从entityPtr中取
*/
func DeleteByCol(colName string, colVal interface{}, statusColName string, entityEmptyPtr TableName) (err error) {
	// auto update cache
	// cols := []string{"status"}

	_, err = envinit.Dbs[common.DBNamespace].Table(entityEmptyPtr).Cols(statusColName).Where(colName+" = ?", colVal).Update(map[string]interface{}{
		statusColName: StatusDeleted, // 此处转小写，否则pg不兼容
	})

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

/**
conditions = ["a = ?"=>1,"b like '%?%'"=>"bb"]
// 如果conditions map的key里不包含“?”，则默认认为value是slice,即使用where in
// 允许conditions里key的value为空
*/
func ListAndCountByCondition(entityPtr interface{}, conditions map[string]interface{}, pageNo int, pageSize int, listSlicePtr interface{}, orderBy string) (total int64, err error) {
	if pageNo <= 0 {
		pageNo = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	whereVal := make([]interface{}, 0)
	whereStr := "true "
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

	session := envinit.Dbs[common.DBNamespace].Where(whereStr, whereVal...)
	for k, v := range whererInMap {
		session.In(k, v)
	}

	var nCount int64
	nCount, err = session.Count(entityPtr)
	if err != nil {
		return
	}

	log.Debugf("nCount of listByCondition =%d", nCount)

	total = nCount

	// 此处必须重新建一个session,不能继续用上次的session
	session = envinit.Dbs[common.DBNamespace].Where(whereStr, whereVal...)
	for k, v := range whererInMap {
		session.In(k, v)
	}
	err = session.OrderBy(orderBy).Limit(pageSize, (pageNo-1)*pageSize).Find(listSlicePtr)
	return
}

/**
conditions = ["a = ?"=>1,"b like '%?%'"=>"bb"]
// 允许conditions里key的value为空
*/
func ListByCondition(entityPtr interface{}, conditions map[string]interface{}, pageNo int, pageSize int, listSlicePtr interface{}, orderBy string) (err error) {

	if pageNo <= 0 {
		pageNo = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	whereVal := make([]interface{}, 0)
	whereStr := "true "
	for ck, cv := range conditions {
		whereStr += " and " + ck
		if cv == nil {
			continue
		}
		whereVal = append(whereVal, cv)
	}

	session := envinit.Dbs[common.DBNamespace].Where(whereStr, whereVal...).Limit(pageSize, (pageNo-1)*pageSize)
	if orderBy != "" {
		session = session.OrderBy(orderBy)
	}

	err = session.Find(listSlicePtr)

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
	whereStr := "true "
	for ck, cv := range conditions {
		whereStr += " and " + ck
		if cv == nil {
			continue
		}
		whereVal = append(whereVal, cv)
	}

	// 每页读取
	for {
		session := envinit.Dbs[common.DBNamespace].Where(whereStr, whereVal...).Limit(pageSize, (startPageNo-1)*pageSize)
		if orderBy != "" {
			session = session.OrderBy(orderBy)
		}

		//listSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(entityPtr)), pageSize, pageSize)
		listSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(tableEntity))).Elem().Addr().Interface()

		err = session.Find(listSlice)
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

/**
conditions = ["a = ?"=>1,"b like '%?%'"=>"bb"]
// 允许conditions里key的value为空
*/
func CountByCondition(entityPtr interface{}, conditions map[string]interface{}) (total int64, err error) {
	whereVal := make([]interface{}, 0)
	whereStr := "true "
	for ck, cv := range conditions {
		whereStr += " and " + ck
		if cv == nil {
			continue
		}
		whereVal = append(whereVal, cv)
	}

	total, err = envinit.Dbs[common.DBNamespace].Where(whereStr, whereVal...).Count(entityPtr)
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

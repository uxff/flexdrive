// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"reflect"
	"strings"
	"time"

	"github.com/go-xorm/core"
)

// Session keep a pointer to sql.DB and provides all execution of all
// kind of database operations.
type Session struct {
	db                     *core.DB
	Engine                 *Engine
	Tx                     *core.Tx
	Statement              Statement
	IsAutoCommit           bool
	IsCommitedOrRollbacked bool
	IsAutoClose            bool

	// Automatically reset the statement after operations that execute a SQL
	// query such as Count(), Find(), Get(), ...
	AutoResetStatement bool

	// !nashtsai! storing these beans due to yet committed tx
	afterInsertBeans map[interface{}]*[]func(interface{})
	afterUpdateBeans map[interface{}]*[]func(interface{})
	afterDeleteBeans map[interface{}]*[]func(interface{})
	// --

	beforeClosures []func(interface{})
	afterClosures  []func(interface{})

	prepareStmt bool
	stmtCache   map[uint32]*core.Stmt //key: hash.Hash32 of (queryStr, len(queryStr))

	// !evalphobia! stored the last executed query on this session
	//beforeSQLExec func(string, ...interface{})
	lastSQL     string
	lastSQLArgs []interface{}

	err error
}

// Clone copy all the session's content and return a new session
func (session *Session) Clone() *Session {
	var sess = *session
	return &sess
}

// Init reset the session as the init status.
func (session *Session) Init() {
	session.Statement.Init()
	session.Statement.Engine = session.Engine
	session.IsAutoCommit = true
	session.IsCommitedOrRollbacked = false
	session.IsAutoClose = false
	session.AutoResetStatement = true
	session.prepareStmt = false

	// !nashtsai! is lazy init better?
	session.afterInsertBeans = make(map[interface{}]*[]func(interface{}), 0)
	session.afterUpdateBeans = make(map[interface{}]*[]func(interface{}), 0)
	session.afterDeleteBeans = make(map[interface{}]*[]func(interface{}), 0)
	session.beforeClosures = make([]func(interface{}), 0)
	session.afterClosures = make([]func(interface{}), 0)

	session.lastSQL = ""
	session.lastSQLArgs = []interface{}{}
}

// Close release the connection from pool
func (session *Session) Close() {
	for _, v := range session.stmtCache {
		v.Close()
	}

	if session.db != nil {
		// When Close be called, if session is a transaction and do not call
		// Commit or Rollback, then call Rollback.
		if session.Tx != nil && !session.IsCommitedOrRollbacked {
			session.Rollback()
		}
		session.Tx = nil
		session.stmtCache = nil
		session.Init()
		session.db = nil
	}
}

func (session *Session) resetStatement() {
	if session.AutoResetStatement {
		session.Statement.Init()
	}
}

// Prepare set a flag to session that should be prepare statement before execute query
func (session *Session) Prepare() *Session {
	session.prepareStmt = true
	return session
}

// Before Apply before Processor, affected bean is passed to closure arg
func (session *Session) Before(closures func(interface{})) *Session {
	if closures != nil {
		session.beforeClosures = append(session.beforeClosures, closures)
	}
	return session
}

// After Apply after Processor, affected bean is passed to closure arg
func (session *Session) After(closures func(interface{})) *Session {
	if closures != nil {
		session.afterClosures = append(session.afterClosures, closures)
	}
	return session
}

// Table can input a string or pointer to struct for special a table to operate.
func (session *Session) Table(tableNameOrBean interface{}) *Session {
	session.Statement.Table(tableNameOrBean)
	return session
}

// Alias set the table alias
func (session *Session) Alias(alias string) *Session {
	session.Statement.Alias(alias)
	return session
}

// NoCascade indicate that no cascade load child object
func (session *Session) NoCascade() *Session {
	session.Statement.UseCascade = false
	return session
}

// ForUpdate Set Read/Write locking for UPDATE
func (session *Session) ForUpdate() *Session {
	session.Statement.IsForUpdate = true
	return session
}

// NoAutoCondition disable generate SQL condition from beans
func (session *Session) NoAutoCondition(no ...bool) *Session {
	session.Statement.NoAutoCondition(no...)
	return session
}

// Limit provide limit and offset query condition
func (session *Session) Limit(limit int, start ...int) *Session {
	session.Statement.Limit(limit, start...)
	return session
}

// OrderBy provide order by query condition, the input parameter is the content
// after order by on a sql statement.
func (session *Session) OrderBy(order string) *Session {
	session.Statement.OrderBy(order)
	return session
}

// Desc provide desc order by query condition, the input parameters are columns.
func (session *Session) Desc(colNames ...string) *Session {
	session.Statement.Desc(colNames...)
	return session
}

// Asc provide asc order by query condition, the input parameters are columns.
func (session *Session) Asc(colNames ...string) *Session {
	session.Statement.Asc(colNames...)
	return session
}

// StoreEngine is only avialble mysql dialect currently
func (session *Session) StoreEngine(storeEngine string) *Session {
	session.Statement.StoreEngine = storeEngine
	return session
}

// Charset is only avialble mysql dialect currently
func (session *Session) Charset(charset string) *Session {
	session.Statement.Charset = charset
	return session
}

// Cascade indicates if loading sub Struct
func (session *Session) Cascade(trueOrFalse ...bool) *Session {
	if len(trueOrFalse) >= 1 {
		session.Statement.UseCascade = trueOrFalse[0]
	}
	return session
}

// NoCache ask this session do not retrieve data from cache system and
// get data from database directly.
func (session *Session) NoCache() *Session {
	session.Statement.UseCache = false
	return session
}

// Join join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (session *Session) Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *Session {
	session.Statement.Join(joinOperator, tablename, condition, args...)
	return session
}

// GroupBy Generate Group By statement
func (session *Session) GroupBy(keys string) *Session {
	session.Statement.GroupBy(keys)
	return session
}

// Having Generate Having statement
func (session *Session) Having(conditions string) *Session {
	session.Statement.Having(conditions)
	return session
}

// DB db return the wrapper of sql.DB
func (session *Session) DB() *core.DB {
	if session.db == nil {
		session.db = session.Engine.db
		session.stmtCache = make(map[uint32]*core.Stmt, 0)
	}
	return session.db
}

func cleanupProcessorsClosures(slices *[]func(interface{})) {
	if len(*slices) > 0 {
		*slices = make([]func(interface{}), 0)
	}
}

func (session *Session) canCache() bool {
	if session.Statement.RefTable == nil ||
		session.Statement.JoinStr != "" ||
		session.Statement.RawSQL != "" ||
		!session.Statement.UseCache ||
		session.Statement.IsForUpdate ||
		session.Tx != nil ||
		len(session.Statement.selectStr) > 0 {
		return false
	}
	return true
}

func (session *Session) doPrepare(sqlStr string) (stmt *core.Stmt, err error) {
	crc := crc32.ChecksumIEEE([]byte(sqlStr))
	// TODO try hash(sqlStr+len(sqlStr))
	var has bool
	stmt, has = session.stmtCache[crc]
	if !has {
		stmt, err = session.DB().Prepare(sqlStr)
		if err != nil {
			return nil, err
		}
		session.stmtCache[crc] = stmt
	}
	return
}

func (session *Session) getField(dataStruct *reflect.Value, key string, table *core.Table, idx int) *reflect.Value {
	var col *core.Column
	if col = table.GetColumnIdx(key, idx); col == nil {
		//session.Engine.logger.Warnf("table %v has no column %v. %v", table.Name, key, table.ColumnsSeq())
		return nil
	}

	fieldValue, err := col.ValueOfV(dataStruct)
	if err != nil {
		session.Engine.logger.Error(err)
		return nil
	}

	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		session.Engine.logger.Warnf("table %v's column %v is not valid or cannot set", table.Name, key)
		return nil
	}
	return fieldValue
}

// Cell cell is a result of one column field
type Cell *interface{}

func (session *Session) rows2Beans(rows *core.Rows, fields []string, fieldsCount int,
	table *core.Table, newElemFunc func([]string) reflect.Value,
	sliceValueSetFunc func(*reflect.Value, core.PK) error) error {
	for rows.Next() {
		var newValue = newElemFunc(fields)
		bean := newValue.Interface()
		dataStruct := rValue(bean)
		// handle beforeClosures
		scanResults, err := session.row2Slice(rows, fields, fieldsCount, bean)
		if err != nil {
			return err
		}
		pk, err := session.slice2Bean(scanResults, fields, fieldsCount, bean, &dataStruct, table)
		if err != nil {
			return err
		}

		err = sliceValueSetFunc(&newValue, pk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (session *Session) row2Slice(rows *core.Rows, fields []string, fieldsCount int, bean interface{}) ([]interface{}, error) {
	for _, closure := range session.beforeClosures {
		closure(bean)
	}

	scanResults := make([]interface{}, fieldsCount)
	for i := 0; i < len(fields); i++ {
		var cell interface{}
		scanResults[i] = &cell
	}
	if err := rows.Scan(scanResults...); err != nil {
		return nil, err
	}

	if b, hasBeforeSet := bean.(BeforeSetProcessor); hasBeforeSet {
		for ii, key := range fields {
			b.BeforeSet(key, Cell(scanResults[ii].(*interface{})))
		}
	}
	return scanResults, nil
}

func (session *Session) slice2Bean(scanResults []interface{}, fields []string, fieldsCount int, bean interface{}, dataStruct *reflect.Value, table *core.Table) (core.PK, error) {
	defer func() {
		if b, hasAfterSet := bean.(AfterSetProcessor); hasAfterSet {
			for ii, key := range fields {
				b.AfterSet(key, Cell(scanResults[ii].(*interface{})))
			}
		}

		// handle afterClosures
		for _, closure := range session.afterClosures {
			closure(bean)
		}
	}()

	var tempMap = make(map[string]int)
	var pk core.PK
	for ii, key := range fields {
		var idx int
		var ok bool
		var lKey = strings.ToLower(key)
		if idx, ok = tempMap[lKey]; !ok {
			idx = 0
		} else {
			idx = idx + 1
		}
		tempMap[lKey] = idx

		if fieldValue := session.getField(dataStruct, key, table, idx); fieldValue != nil {
			rawValue := reflect.Indirect(reflect.ValueOf(scanResults[ii]))

			// if row is null then ignore
			if rawValue.Interface() == nil {
				continue
			}

			if fieldValue.CanAddr() {
				if structConvert, ok := fieldValue.Addr().Interface().(core.Conversion); ok {
					if data, err := value2Bytes(&rawValue); err == nil {
						if err := structConvert.FromDB(data); err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
					continue
				}
			}

			if _, ok := fieldValue.Interface().(core.Conversion); ok {
				if data, err := value2Bytes(&rawValue); err == nil {
					if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
						fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
					}
					fieldValue.Interface().(core.Conversion).FromDB(data)
				} else {
					return nil, err
				}
				continue
			}

			rawValueType := reflect.TypeOf(rawValue.Interface())
			vv := reflect.ValueOf(rawValue.Interface())
			col := table.GetColumnIdx(key, idx)
			if col.IsPrimaryKey {
				pk = append(pk, rawValue.Interface())
			}
			fieldType := fieldValue.Type()
			hasAssigned := false

			if col.SQLType.IsJson() {
				var bs []byte
				if rawValueType.Kind() == reflect.String {
					bs = []byte(vv.String())
				} else if rawValueType.ConvertibleTo(core.BytesType) {
					bs = vv.Bytes()
				} else {
					return nil, fmt.Errorf("unsupported database data type: %s %v", key, rawValueType.Kind())
				}

				hasAssigned = true

				if len(bs) > 0 {
					if fieldValue.CanAddr() {
						err := json.Unmarshal(bs, fieldValue.Addr().Interface())
						if err != nil {
							return nil, err
						}
					} else {
						x := reflect.New(fieldType)
						err := json.Unmarshal(bs, x.Interface())
						if err != nil {
							return nil, err
						}
						fieldValue.Set(x.Elem())
					}
				}

				continue
			}

			switch fieldType.Kind() {
			case reflect.Complex64, reflect.Complex128:
				// TODO: reimplement this
				var bs []byte
				if rawValueType.Kind() == reflect.String {
					bs = []byte(vv.String())
				} else if rawValueType.ConvertibleTo(core.BytesType) {
					bs = vv.Bytes()
				}

				hasAssigned = true
				if len(bs) > 0 {
					if fieldValue.CanAddr() {
						err := json.Unmarshal(bs, fieldValue.Addr().Interface())
						if err != nil {
							return nil, err
						}
					} else {
						x := reflect.New(fieldType)
						err := json.Unmarshal(bs, x.Interface())
						if err != nil {
							return nil, err
						}
						fieldValue.Set(x.Elem())
					}
				}
			case reflect.Slice, reflect.Array:
				switch rawValueType.Kind() {
				case reflect.Slice, reflect.Array:
					switch rawValueType.Elem().Kind() {
					case reflect.Uint8:
						if fieldType.Elem().Kind() == reflect.Uint8 {
							hasAssigned = true
							if col.SQLType.IsText() {
								x := reflect.New(fieldType)
								err := json.Unmarshal(vv.Bytes(), x.Interface())
								if err != nil {
									return nil, err
								}
								fieldValue.Set(x.Elem())
							} else {
								if fieldValue.Len() > 0 {
									for i := 0; i < fieldValue.Len(); i++ {
										if i < vv.Len() {
											fieldValue.Index(i).Set(vv.Index(i))
										}
									}
								} else {
									for i := 0; i < vv.Len(); i++ {
										fieldValue.Set(reflect.Append(*fieldValue, vv.Index(i)))
									}
								}
							}
						}
					}
				}
			case reflect.String:
				if rawValueType.Kind() == reflect.String {
					hasAssigned = true
					fieldValue.SetString(vv.String())
				}
			case reflect.Bool:
				if rawValueType.Kind() == reflect.Bool {
					hasAssigned = true
					fieldValue.SetBool(vv.Bool())
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch rawValueType.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					hasAssigned = true
					fieldValue.SetInt(vv.Int())
				}
			case reflect.Float32, reflect.Float64:
				switch rawValueType.Kind() {
				case reflect.Float32, reflect.Float64:
					hasAssigned = true
					fieldValue.SetFloat(vv.Float())
				}
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				switch rawValueType.Kind() {
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					hasAssigned = true
					fieldValue.SetUint(vv.Uint())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					hasAssigned = true
					fieldValue.SetUint(uint64(vv.Int()))
				}
			case reflect.Struct:
				if fieldType.ConvertibleTo(core.TimeType) {
					dbTZ := session.Engine.DatabaseTZ
					if col.TimeZone != nil {
						dbTZ = col.TimeZone
					}

					if rawValueType == core.TimeType {
						hasAssigned = true

						t := vv.Convert(core.TimeType).Interface().(time.Time)

						z, _ := t.Zone()
						// set new location if database don't save timezone or give an incorrect timezone
						if len(z) == 0 || t.Year() == 0 || t.Location().String() != dbTZ.String() { // !nashtsai! HACK tmp work around for lib/pq doesn't properly time with location
							session.Engine.logger.Debugf("empty zone key[%v] : %v | zone: %v | location: %+v\n", key, t, z, *t.Location())
							t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
								t.Minute(), t.Second(), t.Nanosecond(), dbTZ)
						}

						t = t.In(session.Engine.TZLocation)
						fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
					} else if rawValueType == core.IntType || rawValueType == core.Int64Type ||
						rawValueType == core.Int32Type {
						hasAssigned = true

						t := time.Unix(vv.Int(), 0).In(session.Engine.TZLocation)
						fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
					} else {
						if d, ok := vv.Interface().([]uint8); ok {
							hasAssigned = true
							t, err := session.byte2Time(col, d)
							if err != nil {
								session.Engine.logger.Error("byte2Time error:", err.Error())
								hasAssigned = false
							} else {
								fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
							}
						} else if d, ok := vv.Interface().(string); ok {
							hasAssigned = true
							t, err := session.str2Time(col, d)
							if err != nil {
								session.Engine.logger.Error("byte2Time error:", err.Error())
								hasAssigned = false
							} else {
								fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
							}
						} else {
							return nil, fmt.Errorf("rawValueType is %v, value is %v", rawValueType, vv.Interface())
						}
					}
				} else if nulVal, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
					// !<winxxp>! 增加支持sql.Scanner接口的结构，如sql.NullString
					hasAssigned = true
					if err := nulVal.Scan(vv.Interface()); err != nil {
						session.Engine.logger.Error("sql.Sanner error:", err.Error())
						hasAssigned = false
					}
				} else if col.SQLType.IsJson() {
					if rawValueType.Kind() == reflect.String {
						hasAssigned = true
						x := reflect.New(fieldType)
						if len([]byte(vv.String())) > 0 {
							err := json.Unmarshal([]byte(vv.String()), x.Interface())
							if err != nil {
								return nil, err
							}
							fieldValue.Set(x.Elem())
						}
					} else if rawValueType.Kind() == reflect.Slice {
						hasAssigned = true
						x := reflect.New(fieldType)
						if len(vv.Bytes()) > 0 {
							err := json.Unmarshal(vv.Bytes(), x.Interface())
							if err != nil {
								return nil, err
							}
							fieldValue.Set(x.Elem())
						}
					}
				} else if session.Statement.UseCascade {
					table, err := session.Engine.autoMapType(*fieldValue)
					if err != nil {
						return nil, err
					}

					hasAssigned = true
					if len(table.PrimaryKeys) != 1 {
						return nil, errors.New("unsupported non or composited primary key cascade")
					}
					var pk = make(core.PK, len(table.PrimaryKeys))
					pk[0], err = asKind(vv, rawValueType)
					if err != nil {
						return nil, err
					}

					if !isPKZero(pk) {
						// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
						// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
						// property to be fetched lazily
						structInter := reflect.New(fieldValue.Type())
						newsession := session.Engine.NewSession()
						defer newsession.Close()
						has, err := newsession.ID(pk).NoCascade().Get(structInter.Interface())
						if err != nil {
							return nil, err
						}
						if has {
							fieldValue.Set(structInter.Elem())
						} else {
							return nil, errors.New("cascade obj is not exist")
						}
					}
				}
			case reflect.Ptr:
				// !nashtsai! TODO merge duplicated codes above
				switch fieldType {
				// following types case matching ptr's native type, therefore assign ptr directly
				case core.PtrStringType:
					if rawValueType.Kind() == reflect.String {
						x := vv.String()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrBoolType:
					if rawValueType.Kind() == reflect.Bool {
						x := vv.Bool()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrTimeType:
					if rawValueType == core.PtrTimeType {
						hasAssigned = true
						var x = rawValue.Interface().(time.Time)
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrFloat64Type:
					if rawValueType.Kind() == reflect.Float64 {
						x := vv.Float()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUint64Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint64(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt64Type:
					if rawValueType.Kind() == reflect.Int64 {
						x := vv.Int()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrFloat32Type:
					if rawValueType.Kind() == reflect.Float64 {
						var x = float32(vv.Float())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrIntType:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUintType:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUint32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Complex64Type:
					var x complex64
					if len([]byte(vv.String())) > 0 {
						err := json.Unmarshal([]byte(vv.String()), &x)
						if err != nil {
							return nil, err
						}
						fieldValue.Set(reflect.ValueOf(&x))
					}
					hasAssigned = true
				case core.Complex128Type:
					var x complex128
					if len([]byte(vv.String())) > 0 {
						err := json.Unmarshal([]byte(vv.String()), &x)
						if err != nil {
							return nil, err
						}
						fieldValue.Set(reflect.ValueOf(&x))
					}
					hasAssigned = true
				} // switch fieldType
			} // switch fieldType.Kind()

			// !nashtsai! for value can't be assigned directly fallback to convert to []byte then back to value
			if !hasAssigned {
				data, err := value2Bytes(&rawValue)
				if err != nil {
					return nil, err
				}

				if err = session.bytes2Value(col, fieldValue, data); err != nil {
					return nil, err
				}
			}
		}
	}
	return pk, nil
}

func (session *Session) queryPreprocess(sqlStr *string, paramStr ...interface{}) {
	for _, filter := range session.Engine.dialect.Filters() {
		*sqlStr = filter.Do(*sqlStr, session.Engine.dialect, session.Statement.RefTable)
	}

	session.saveLastSQL(*sqlStr, paramStr...)
}

// saveLastSQL stores executed query information
func (session *Session) saveLastSQL(sql string, args ...interface{}) {
	session.lastSQL = sql
	session.lastSQLArgs = args
	session.Engine.logSQL(sql, args...)
}

// LastSQL returns last query information
func (session *Session) LastSQL() (string, []interface{}) {
	return session.lastSQL, session.lastSQLArgs
}

// tbName get some table's table name
func (session *Session) tbNameNoSchema(table *core.Table) string {
	if len(session.Statement.AltTableName) > 0 {
		return session.Statement.AltTableName
	}

	return table.Name
}

// Unscoped always disable struct tag "deleted"
func (session *Session) Unscoped() *Session {
	session.Statement.Unscoped()
	return session
}

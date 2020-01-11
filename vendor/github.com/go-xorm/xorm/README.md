[中文](https://github.com/go-xorm/xorm/blob/master/README_CN.md)

Xorm is a simple and powerful ORM for Go.

[![CircleCI](https://circleci.com/gh/go-xorm/xorm.svg?style=shield)](https://circleci.com/gh/go-xorm/xorm) [![codecov](https://codecov.io/gh/go-xorm/xorm/branch/master/graph/badge.svg)](https://codecov.io/gh/go-xorm/xorm)
[![](https://goreportcard.com/badge/github.com/go-xorm/xorm)](https://goreportcard.com/report/github.com/go-xorm/xorm) 
[![Join the chat at https://img.shields.io/discord/323460943201959939.svg](https://img.shields.io/discord/323460943201959939.svg)](https://discord.gg/HuR2CF3)

# Notice

The last master version is not backwards compatible. You should use `engine.ShowSQL()` and `engine.Logger().SetLevel()` instead of `engine.ShowSQL = `, `engine.ShowInfo = ` and so on.

# Features

* Struct <-> Table Mapping Support

* Chainable APIs

* Transaction Support

* Both ORM and raw SQL operation Support

* Sync database schema Support

* Query Cache speed up

* Database Reverse support, See [Xorm Tool README](https://github.com/go-xorm/cmd/blob/master/README.md)

* Simple cascade loading support

* Optimistic Locking support

* SQL Builder support via [github.com/go-xorm/builder](https://github.com/go-xorm/builder)

# Drivers Support

Drivers for Go's sql package which currently support database/sql includes:

* Mysql: [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)

* MyMysql: [github.com/ziutek/mymysql/godrv](https://github.com/ziutek/mymysql/godrv)

* Postgres: [github.com/lib/pq](https://github.com/lib/pq)

* Tidb: [github.com/pingcap/tidb](https://github.com/pingcap/tidb)

* SQLite: [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

* MsSql: [github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb)

* Oracle: [github.com/mattn/go-oci8](https://github.com/mattn/go-oci8) (experiment)

# Changelog

* **v0.6.2**
    * refactor tag parse methods
    * add Scan features to Get
    * add QueryString method

* **v0.6.0**
    * remove support for ql
    * add query condition builder support via [github.com/go-xorm/builder](https://github.com/go-xorm/builder), so `Where`, `And`, `Or` 
methods can use `builder.Cond` as parameter
    * add Sum, SumInt, SumInt64 and NotIn methods
    * some bugs fixed

* **v0.5.0**
    * logging interface changed
    * some bugs fixed

* **v0.4.5**
    * many bugs fixed
    * extends support unlimited deepth
    * Delete Limit support

* **v0.4.4**
    * ql database expriment support
    * tidb database expriment support
    * sql.NullString and etc. field support
    * select ForUpdate support
    * many bugs fixed

[More changes ...](https://github.com/go-xorm/manual-en-US/tree/master/chapter-16)

# Installation

	go get github.com/go-xorm/xorm

# Documents

* [Manual](http://xorm.io/docs)

* [GoDoc](http://godoc.org/github.com/go-xorm/xorm)

* [GoWalker](http://gowalker.org/github.com/go-xorm/xorm)

# Quick Start

* Create Engine

```Go
engine, err := xorm.NewEngine(driverName, dataSourceName)
```

* Define a struct and Sync2 table struct to database

```Go
type User struct {
    Id int64
    Name string
    Salt string
    Age int
    Passwd string `xorm:"varchar(200)"`
    Created time.Time `xorm:"created"`
    Updated time.Time `xorm:"updated"`
}

err := engine.Sync2(new(User))
```

* `Query` runs a SQL string, the returned results is `[]map[string][]byte`, `QueryString` returns `[]map[string]string`.

```Go
results, err := engine.Query("select * from user")

results, err := engine.QueryString("select * from user")
```

* `Execute` runs a SQL string, it returns `affetcted` and `error`

```Go
affected, err := engine.Exec("update user set age = ? where name = ?", age, name)
```

* `Insert` one or multiple records to database

```Go
affected, err := engine.Insert(&user)
// INSERT INTO struct () values ()
affected, err := engine.Insert(&user1, &user2)
// INSERT INTO struct1 () values ()
// INSERT INTO struct2 () values ()
affected, err := engine.Insert(&users)
// INSERT INTO struct () values (),(),()
affected, err := engine.Insert(&user1, &users)
// INSERT INTO struct1 () values ()
// INSERT INTO struct2 () values (),(),()
```

* Query one record from database

```Go
has, err := engine.Get(&user)
// SELECT * FROM user LIMIT 1
has, err := engine.Where("name = ?", name).Desc("id").Get(&user)
// SELECT * FROM user WHERE name = ? ORDER BY id DESC LIMIT 1
var name string
has, err := engine.Where("id = ?", id).Cols("name").Get(&name)
// SELECT name FROM user WHERE id = ?
var id int64
has, err := engine.Where("name = ?", name).Cols("id").Get(&id)
// SELECT id FROM user WHERE name = ?
var valuesMap = make(map[string]string)
has, err := engine.Where("id = ?", id).Get(&valuesMap)
// SELECT * FROM user WHERE id = ?
var valuesSlice = make([]interface{}, len(cols))
has, err := engine.Where("id = ?", id).Cols(cols...).Get(&valuesSlice)
// SELECT col1, col2, col3 FROM user WHERE id = ?
```

* Query multiple records from database, also you can use join and extends

```Go
var users []User
err := engine.Where("name = ?", name).And("age > 10").Limit(10, 0).Find(&users)
// SELECT * FROM user WHERE name = ? AND age > 10 limit 0 offset 10

type Detail struct {
    Id int64
    UserId int64 `xorm:"index"`
}

type UserDetail struct {
    User `xorm:"extends"`
    Detail `xorm:"extends"`
}

var users []UserDetail
err := engine.Table("user").Select("user.*, detail.*")
    Join("INNER", "detail", "detail.user_id = user.id").
    Where("user.name = ?", name).Limit(10, 0).
    Find(&users)
// SELECT user.*, detail.* FROM user INNER JOIN detail WHERE user.name = ? limit 0 offset 10
```

* Query multiple records and record by record handle, there are two methods Iterate and Rows

```Go
err := engine.Iterate(&User{Name:name}, func(idx int, bean interface{}) error {
    user := bean.(*User)
    return nil
})
// SELECT * FROM user

rows, err := engine.Rows(&User{Name:name})
// SELECT * FROM user
defer rows.Close()
bean := new(Struct)
for rows.Next() {
    err = rows.Scan(bean)
}
```

* Update one or more records, default will update non-empty and non-zero fields except when you use Cols, AllCols and so on.

```Go
affected, err := engine.Id(1).Update(&user)
// UPDATE user SET ... Where id = ?

affected, err := engine.Update(&user, &User{Name:name})
// UPDATE user SET ... Where name = ?

var ids = []int64{1, 2, 3}
affected, err := engine.In("id", ids).Update(&user)
// UPDATE user SET ... Where id IN (?, ?, ?)

// force update indicated columns by Cols
affected, err := engine.Id(1).Cols("age").Update(&User{Name:name, Age: 12})
// UPDATE user SET age = ?, updated=? Where id = ?

// force NOT update indicated columns by Omit
affected, err := engine.Id(1).Omit("name").Update(&User{Name:name, Age: 12})
// UPDATE user SET age = ?, updated=? Where id = ?

affected, err := engine.Id(1).AllCols().Update(&user)
// UPDATE user SET name=?,age=?,salt=?,passwd=?,updated=? Where id = ?
```

* Delete one or more records, Delete MUST have condition

```Go
affected, err := engine.Where(...).Delete(&user)
// DELETE FROM user Where ...
affected, err := engine.Id(2).Delete(&user)
```

* Count records

```Go
counts, err := engine.Count(&user)
// SELECT count(*) AS total FROM user
```

* Query conditions builder

```Go
err := engine.Where(builder.NotIn("a", 1, 2).And(builder.In("b", "c", "d", "e"))).Find(&users)
// SELECT id, name ... FROM user WHERE a NOT IN (?, ?) AND b IN (?, ?, ?)
```

# Cases

* [github.com/m3ng9i/qreader](https://github.com/m3ng9i/qreader)

* [Wego](http://github.com/go-tango/wego)

* [Docker.cn](https://docker.cn/)

* [Gogs](http://try.gogits.org) - [github.com/gogits/gogs](http://github.com/gogits/gogs)

* [Gorevel](http://gorevel.cn/) - [github.com/goofcc/gorevel](http://github.com/goofcc/gorevel)

* [Gowalker](http://gowalker.org) - [github.com/Unknwon/gowalker](http://github.com/Unknwon/gowalker)

* [Gobuild.io](http://gobuild.io) - [github.com/shxsun/gobuild](http://github.com/shxsun/gobuild)

* [Sudo China](http://sudochina.com) - [github.com/insionng/toropress](http://github.com/insionng/toropress)

* [Godaily](http://godaily.org) - [github.com/govc/godaily](http://github.com/govc/godaily)

* [YouGam](http://www.yougam.com/)

* [GoCMS - github.com/zzboy/GoCMS](https://github.com/zzdboy/GoCMS)

* [GoBBS - gobbs.domolo.com](http://gobbs.domolo.com/)

* [go-blog](http://wangcheng.me) - [github.com/easykoo/go-blog](https://github.com/easykoo/go-blog)

# Discuss

Please visit [Xorm on Google Groups](https://groups.google.com/forum/#!forum/xorm)

# Contributing

If you want to pull request, please see [CONTRIBUTING](https://github.com/go-xorm/xorm/blob/master/CONTRIBUTING.md)

# LICENSE

 BSD License
 [http://creativecommons.org/licenses/BSD/](http://creativecommons.org/licenses/BSD/)

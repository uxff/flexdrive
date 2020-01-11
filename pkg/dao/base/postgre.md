# postgre 迁移踩坑

### 1、where 子句 
where 1 and mStatus=1
要改成
where true and mStatus=1

### 2、类型和映射
不要开启其他column mapper的时候：
字段名数据库中是agId 则代码entity定义的名称必须是AgId 除了首字母，其他都要和表中完全一致。

如果开启了则按照mapper走。
xorm定义的 tag 中 '字段名1' 这种生效，在pg中会被加引号。

### 3、时间不能有0000-00-00
xorm定义的default('0001-01-01 00:00:00') 不生效。
写入时，判断如果time.Unix()<1000则time.Add(1*time.Second)

### 4、并行写太多会丢数据，报too many clients
show max_connections;返回2100
测试10个协程不会有问题，100个协程会丢数据。


### 5、最终方案
5.1 数据库ddl字段名使用小写。
entity定义的变量名不变。

5.2 entity声明增加tag tag中字段名使用小写。
todo: 测试mysql兼容性

5.3 开发中的字段不要加引号，如果一定要加，那就加成小写。

5.4 不要使用columeMapper。

postgres 错误duplicate key value violates unique constraint 解决方案
SELECT setval('tablename_id_seq', (SELECT MAX(id) FROM tablename)+1)


主要是：serial key其实是由sequence实现的，当你手动给serial列赋值的时候，sequence是不会自增量变化的。
最好不要给serial手工赋值

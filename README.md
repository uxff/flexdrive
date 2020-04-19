# flexdrive
一个分布式云盘系统。

后端采用golang的gin框架，前端使用Bootstrap，集群通信使用gRPC协议。

具有以下特性：
- a. 一个节点同时提供web和gRPC存储服务；
- b. 可无限水平扩展，通过在分布式集群增加存储节点的方式实现扩展，理论上空间无限制；
- c. 一个文件被随机分布在至少三个节点上，避免单点故障导致文件丢失，增加服务可靠性；
- d. 集群采用分布式负载均衡，将文件上传和备份到负载最小的节点上；
- e. 集群使用Golang作为主要开发语言，编译后无依赖，可容器化部署运行；
- f. 集群间使用gRPC通信协议，用于心跳，同步集群节点信息等；
- g. 提供后台管理界面，满足内容可控需求。


集群运行结构图:
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture.png">
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture2.png">
节点内部结构图；
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture3.png">

数据库支持mysql，可使用sqlite替代。

## 开始
准备好环境变量
|--:|--:|--:|
| 环境变量	|示例值	|含义说明|
|--:|--:|--:|
|SERVEADMIN	|127.0.0.1:10011	|管理端监听端口|
|SERVECUSTOMER	|127.0.0.1:10012	|会员端监听端口|
|SERVECLUSTER	|127.0.0.1:10013	|节点通信监听端口|
|DATADSN	|mysql://user:pass@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4	|关系数据库的DSN|
|STORAGEDIR	|./data/	|物理文件存储目录|

编译：
```
$ go run cmd/node/main.go
```

单节点运行：
```
$ SERVEADMIN=127.0.0.1:10011 SERVECUSTOMER=127.0.0.1:10012 SERVECLUSTER=127.0.0.1:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./main
```


多节点运行；
```
sh runcluster.sh
```
启动后访问管理端 http://127.0.0.1:10011 访问会员端 http://127.0.0.1:10012


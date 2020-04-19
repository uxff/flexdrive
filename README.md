# flexdrive
一个分布式云盘系统。

后端采用golang的gin框架，前端使用Bootstrap，集群通信使用gRPC协议。

具有以下特性：
a. 一个节点同时提供web和gRPC存储服务；
b. 可无限水平扩展，通过在分布式集群增加存储节点的方式实现扩展，理论上空间无限制；
c. 一个文件被随机分布在至少三个节点上，避免单点故障导致文件丢失，增加服务可靠性；
d. 集群采用分布式负载均衡，将文件上传和备份到负载最小的节点上；
e. 集群使用Golang作为主要开发语言，编译后无依赖，可容器化部署运行；
f. 集群间使用gRPC通信协议，用于心跳，同步集群节点信息等；
g. 提供后台管理界面，满足内容可控需求。


集群运行结构图:
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture.png">
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture2.png">
节点内部结构图；
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture3.png">




# Flexdrive
A distributed cloud disk system.

golang gin framework is used in the back end, Bootstrap is used in the front end, and gRPC protocol is used for cluster communication.

Has the following characteristics:

- a. One node provides both web and gRPC storage services.
- b. It can be expanded horizontally indefinitely by adding storage nodes in a distributed cluster. Theoretically, the space is unlimited.
- c. A file is randomly distributed on at least three nodes to prevent file loss caused by single point of failure and improve service reliability.
- d. The cluster uses distributed load balancing to upload and back up files to nodes with the least load.
- e. The cluster uses Golang as the main development language, which is independent after compilation and can be deployed and run in container;
- f. The gRPC communication protocol is used between clusters for heartbeat and node synchronization.
- g. Provide a background management interface to meet content controllable requirements.


Cluster operation structure diagram:
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture.png">
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture2.png">
Node internal structure diagram:
<img src="https://github.com/uxff/flexdrive/raw/master/static/images/clusters-architecture3.png">

It supports mysql for database, you can use sqlite instead.

## Get Started
Prepare your environment variables

|--:|--:|--:|
| Env Variables	|Example	|Description|
|--:|--:|--:|
|SERVEADMIN	|127.0.0.1:10011	|address  and port for administrators|
|SERVECUSTOMER	|127.0.0.1:10012	|address and port for guests|
|SERVECLUSTER	|127.0.0.1:10013	|address and port for cluster communication|
|DATADSN	|mysql://user:pass@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4	|DNS of database|
|STORAGEDIR	|./data/	|physical file storage|

Compiling:
```
$ go run cmd/node/main.go
```

Run by single node:
```
$ SERVEADMIN=127.0.0.1:10011 SERVECUSTOMER=127.0.0.1:10012 SERVECLUSTER=127.0.0.1:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./main
```


Run by multing nodes:
```
sh runcluster.sh
```

Visit admin web via: http://127.0.0.1:10011 

Visit guest web via: http://127.0.0.1:10012


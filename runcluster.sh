#go build -o flexdrive cmd/node/main.go
SERVEWEB=0.0.0.0:10011 SERVECLUSTER=0.0.0.0:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./flexdrive > 1.log 2>&1 &
SERVEWEB=0.0.0.0:10021 SERVECLUSTER=0.0.0.0:10023 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./flexdrive > 2.log 2>&1 &
SERVEWEB=0.0.0.0:10031 SERVECLUSTER=0.0.0.0:10033 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.1:10023,127.0.0.1:10033 DATADSN='sqlite3://./flexdrive.db'  STORAGEDIR=./data/ ./flexdrive > 3.log 2>&1 &

ps aux | grep flexdrive | grep -v grep
echo Visit admin:    http://127.0.0.1:10011/adm
echo Visit customer: http://127.0.0.1:10012/

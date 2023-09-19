go build -o flexdrive cmd/node/main.go
GIN_MODE=release SERVEADMIN=127.0.0.1:10011 SERVECUSTOMER=127.0.0.1:10012 SERVECLUSTER=127.0.0.1:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.2:10013,127.0.0.3:10013 DATADSN='mysql://root:123x456@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4&parseTime=True&loc=Local'  STORAGEDIR=./data/ ./flexdrive > 1.log 2>&1 &
GIN_MODE=release SERVEADMIN=127.0.0.2:10011 SERVECUSTOMER=127.0.0.2:10012 SERVECLUSTER=127.0.0.2:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.2:10013,127.0.0.3:10013 DATADSN='mysql://root:123x456@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4&parseTime=True&loc=Local'  STORAGEDIR=./data/ ./flexdrive > 2.log 2>&1 &
GIN_MODE=release SERVEADMIN=127.0.0.3:10011 SERVECUSTOMER=127.0.0.3:10012 SERVECLUSTER=127.0.0.3:10013 CLUSTERMEMBERS=127.0.0.1:10013,127.0.0.2:10013,127.0.0.3:10013 DATADSN='mysql://root:123x456@tcp(127.0.0.1:3306)/flexdrive?charset=utf8mb4&parseTime=True&loc=Local'  STORAGEDIR=./data/ ./flexdrive > 3.log 2>&1 &


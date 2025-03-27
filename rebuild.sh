ps aux | grep flexdrive | grep -v grep
echo kill ` ps x | grep flexdrive | grep -v grep | awk '{print $1}'`
kill ` ps x | grep flexdrive | grep -v grep | awk '{print $1}'`
ps aux | grep flexdrive | grep -v grep

git pull
go build -o flexdrive cmd/node/main.go
sh ./runcluster-mysql.sh 
#sleep 1
#execute below in windows will cause flexdrive terminate. do not run below.
#tail -f 1.log

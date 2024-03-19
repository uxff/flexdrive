ps aux | grep flexdrive | grep -v grep
kill ` ps x | grep flexdrive | grep -v grep | awk '{print $1}'`
ps aux | grep flexdrive | grep -v grep

git pull
go build -o flexdrive cmd/node/main.go
sh ./runcluster.sh
tail -f 1.log

ps aux | grep flexdrive
kill ` ps aux | grep flexdrive | grep -v grep | awk '{print $2}'`
ps aux | grep flexdrive

git pull
go build -o flexdrive cmd/node/main.go
sh ./runcluster.sh
tail -f 1.log

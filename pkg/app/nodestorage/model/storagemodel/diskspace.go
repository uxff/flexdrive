package storagemodel

import (
	"fmt"
	"os"
)

// 不能实时获取 应该命令行启动时候传入
// type DiskStatus struct {
// 	All  uint64 `json:"all"`
// 	Used uint64 `json:"used"`
// 	Free uint64 `json:"free"`
// }

// // disk usage of path/disk
// func DiskUsage(path string) (disk DiskStatus) {
// 	fs := syscall.Statfs_t{}
// 	err := syscall.Statfs(path, &fs)
// 	if err != nil {
// 		return
// 	}
// 	disk.All = fs.Blocks * uint64(fs.Bsize)
// 	disk.Free = fs.Bfree * uint64(fs.Bsize)
// 	disk.Used = disk.All - disk.Free
// 	return

// }

func DirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

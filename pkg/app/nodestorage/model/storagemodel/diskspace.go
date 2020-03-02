package storagemodel

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/uxff/flexdrive/pkg/log"
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
		log.Debugf("%v", err)
		return false
	}
	return true
}

// func downloadFile(url string, localPath string, fb func(length, downLen int64)) error {
// 	    var (
// 	        fsize   int64
// 	        buf     = make([]byte, 32*1024)
// 	        written int64
// 	    )
// 	    tmpFilePath := localPath + ".download"
// 	    log.Debugf("%v",tmpFilePath)
// 	    //创建一个http client
// 	    client := new(http.Client)
// 	    //client.Timeout = time.Second * 60 //设置超时时间
// 	    //get方法获取资源
// 	    resp, err := client.Get(url)
// 	    if err != nil {
// 	        return err
// 	    }

// 	    //读取服务器返回的文件大小
// 	    fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
// 	    if err != nil {
// 	        log.Debugf("%v",err)
// 	    }
// 	    if IsFileExist(localPath, fsize) {
// 	        return err
// 	    }
// 	    log.Debugf("fsize: %v", fsize)
// 	    //创建文件
// 	    file, err := os.Create(tmpFilePath)
// 	    if err != nil {
// 	        return err
// 	    }
// 	    defer file.Close()
// 	    if resp.Body == nil {
// 	        return errors.New("body is null")
// 	    }
// 	    defer resp.Body.Close()
// 	    //下面是 io.copyBuffer() 的简化版本
// 	    for {
// 	        //读取bytes
// 	        nr, er := resp.Body.Read(buf)
// 	        if nr > 0 {
// 	            //写入bytes
// 	            nw, ew := file.Write(buf[0:nr])
// 	            //数据长度大于0
// 	            if nw > 0 {
// 	                written += int64(nw)
// 	            }
// 	            //写入出错
// 	            if ew != nil {
// 	                err = ew
// 	                break
// 	            }
// 	            //读取是数据长度不等于写入的数据长度
// 	            if nr != nw {
// 	                err = io.ErrShortWrite
// 	                break
// 	            }
// 	        }
// 	        if er != nil {
// 	            if er != io.EOF {
// 	                err = er
// 	            }
// 	            break
// 	        }
// 	        //没有错误了快使用 callback
// 	        fb(fsize, written)
// 	    }
// 	    log.Debugf("%v",err)
// 	    if err == nil {
// 	        file.Close()
// 	        err = os.Rename(tmpFilePath, localPath)
// 	        log.Debugf("%v",err)
// 	    }
// 	    return err
// }

func downloadFile(fileUrl string, localPath string) error {
	// nt := time.Now().Format("2006-01-02 15:04:05")
	// fmt.Printf("[%s]To download %s\n", nt, fileID)

	//url := fmt.Sprintf("http://%s/file/%s", node, fileID)
	// localPath := fmt.Sprintf("/yourpath/download/%s_%s", clustername, fileID)
	newFile, err := os.Create(localPath)
	if err != nil {
		log.Errorf("create localFile %s error: %v", localPath, err)
		return err
	}
	defer newFile.Close()

	client := http.Client{Timeout: 900 * time.Second}
	resp, err := client.Get(fileUrl)
	defer resp.Body.Close()

	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		log.Errorf("when downloadFile %s error:%v", fileUrl, err)
		return err
	}
	return nil
}

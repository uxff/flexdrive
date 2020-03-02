package storagemodel

import (
	"io"
	"os"
	"strings"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/httpworker"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/filehash"

	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
)

const (
	DirSplitDeep = 2 // 文件切割深度
)

const (
	DefaultStorageDir = "/tmp/flexdrive/"
)

type NodeStorage struct {
	StorageDir     string // 本节点的存储路径 保证有/结尾
	ClusterId      string
	ClusterMembers string
	WorkerAddr     string

	NodeEnt *dao.Node
	Worker  *httpworker.Worker
}

var node *NodeStorage

func init() {
	node = &NodeStorage{}
}

//
func StartNode(storageDir string, httpAddr string, clusterId string, clusterMembers string) error {
	if storageDir == "" {
		storageDir = DefaultStorageDir
	}

	if storageDir[len(storageDir)-1] != '/' {
		storageDir += "/"
	}

	node.NodeEnt = &dao.Node{}

	node.StorageDir = storageDir
	node.ClusterId = clusterId
	node.ClusterMembers = clusterMembers
	node.WorkerAddr = httpAddr

	// 准备makedir

	if !DirExist(node.StorageDir) {
		err := os.MkdirAll(node.StorageDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}

	node.Worker = httpworker.NewWorker(node.WorkerAddr, node.ClusterId)
	node.Worker.AddMates(strings.Split(node.ClusterMembers, ","))
	node.NodeEnt.NodeName = node.Worker.Id

	// 准备启动服务
	serveErrorChan := make(chan error, 1)

	// start http server
	go func() {
		log.Debugf("http server will start at %v", node.WorkerAddr)
		serveErrorChan <- node.Worker.ServePingable()
	}()

	// start cluster node
	go func() {
		log.Debugf("worker server will start ")
		serveErrorChan <- node.Worker.Start()
	}()

	err := <-serveErrorChan
	log.Errorf("an error occur when serving storage: %v", err)

	return err

	// 监听信号，先关闭rpc服务，再关闭消息队列
	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	// select {
	// case sig := <-ch:
	// 	log.Debugf("receive signal '%v', server will exit", sig)
	// 	node.Worker.Quit()
	// }
	// log.Debugf("start node, storageDir=%s", node.StorageDir)

	// return nil
}

// func Register() {

// }

// func KeepRegistered() {

// }

// func Accept() {

// }

// 获取节点可用空闲空间
func (n *NodeStorage) GetFreeSpace() int64 {
	return 1024 * 1024 * 1024
}

// 将文件保存到本地 fileHash用于校验
func (n *NodeStorage) SaveFile(filepath string, fileHash string) (*dao.FileIndex, error) {
	fileHandle, err := os.Open(filepath)
	if err != nil {
		log.Errorf("open %s failed:%v", filepath, err)
		return nil, err
	}

	defer fileHandle.Close()

	fileHash, err = filehash.CalcFileSha1(fileHandle)
	if err != nil {
		log.Errorf("calc filehash if uploaded file failed:%v", err)
		//StdErrResponse(c, ErrInternal)
		return nil, err
	}

	return &dao.FileIndex{
		FileName: filepath,
		FileHash: fileHash,
	}, nil

}

func (n *NodeStorage) SaveFileFromNode(filepath string, nodeAddr string) (*dao.FileIndex, error) {

	return &dao.FileIndex{
		FileName: filepath,
		//FileHash: fileHash,
	}, nil
}

// 保存已经打开的文件流 必须在调用前知道fileHash
// fileName 可空
func (n *NodeStorage) SaveFileHandler(inputFileHandler io.Reader, fileHash string, fileName string, size int64) (*dao.FileIndex, error) {

	fileInStorage := n.FileHashToStoragePath(fileHash) //n.StorageDir + fileHash

	// 将inputFileHandler 的保存在节点存储系统中
	fileInStorageHandle, err := os.OpenFile(fileInStorage, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		// 不能打开临时分拣
		log.Errorf("open %s failed:%v", fileInStorage, err)
		return nil, err
	}

	defer fileInStorageHandle.Close()

	_, err = io.Copy(fileInStorageHandle, inputFileHandler)
	if err != nil {
		log.Errorf("file copy from uploaded tmp handle to %s failed:%v", fileInStorage, err)
		return nil, err
	}

	// fileInStorageHandle.Seek(0, io.SeekStart)
	// fileHash, err := filehash.CalcFileSha1(headerFileHandle)
	// if err != nil {
	// 	log.Trace(requestId).Errorf("calc filehash if uploaded file failed:%v", err)
	// 	StdErrResponse(c, ErrInternal)
	// 	return
	// }

	fileIndex := &dao.FileIndex{
		FileName:  fileName,
		FileHash:  fileHash,
		NodeId:    node.NodeEnt.Id,
		InnerPath: fileInStorage,
		OuterPath: "", // todo
		Size:      size,
		Space:     size / 1024,
	}

	_, err = base.Insert(fileIndex)
	if err != nil {
		log.Errorf("insert %s failed:%v", fileInStorage, err)
		return nil, err
	}

	log.Infof("a file stored, %s", fileInStorage)
	// todo distribute to other node

	return fileIndex, nil
}

func GetCurrentNode() *NodeStorage {
	return node
}

func (n *NodeStorage) FileHashToStoragePath(fileHash string) string {
	//return node.StorageDir + fileHash

	splitDeep := 2
	curDir := node.StorageDir

	for deepIdx := 0; deepIdx < splitDeep; deepIdx++ {

		prefix1 := fileHash[deepIdx : deepIdx+1]
		//prefix2 := fileHash[1:2]
		curDir = curDir + prefix1 + "/"
		if !DirExist(curDir) {
			os.MkdirAll(curDir, os.ModeDir|os.ModePerm)
		}
	}

	return curDir + fileHash

}

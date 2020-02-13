package storagemodel

import (
	"io"
	"os"

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
	NodeEnt    *dao.Node
	StorageDir string // 本节点的存储路径 保证有/结尾
}

var node *NodeStorage

func init() {
	node = &NodeStorage{}
}

//
func StartNode(name string, storageDir string) {
	if storageDir == "" {
		storageDir = DefaultStorageDir
	}

	if storageDir[len(storageDir)-1] != '/' {
		storageDir += "/"
	}

	node.NodeEnt = &dao.Node{
		NodeName: name,
	}

	node.StorageDir = storageDir

	// 准备makedir

	if !DirExist(node.StorageDir) {
		os.MkdirAll(node.StorageDir, os.ModeDir)
	}
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

func (n *NodeStorage) SaveFile(filepath string, fileHash string) (*dao.FileIndex, error) {
	return &dao.FileIndex{
		FileName: filepath,
		FileHash: fileHash,
	}, nil

}

// fileName 可空
func (n *NodeStorage) SaveFileHandler(inputFileHandler io.Reader, fileHash string, fileName string, size int64) (*dao.FileIndex, error) {

	fileInStorage := "/tmp/flexdrive/" + fileHash

	// 将inputFileHandler 的保存在节点存储系统中
	fileInStorageHandle, err := os.Open(fileInStorage)
	if err != nil {
		// 不能打开临时分拣
		return nil, err
	}

	defer fileInStorageHandle.Close()

	_, err = io.Copy(fileInStorageHandle, inputFileHandler)

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
		return nil, err
	}

	// todo distribute to other node

	return fileIndex, nil
}

func GetCurrentNode() *NodeStorage {
	return &NodeStorage{}
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
			os.MkdirAll(curDir, os.ModeDir)
		}
	}

	return curDir + fileHash

}

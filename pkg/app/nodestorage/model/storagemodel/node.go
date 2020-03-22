package storagemodel

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	worker "github.com/uxff/flexdrive/pkg/app/nodestorage/httpworker"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/filehash"

	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
)

const (
	DirSplitDeep = 2 // 文件存储目录深度
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
	Worker  *worker.Worker
}

var node *NodeStorage

func init() {
	node = &NodeStorage{}
}

// 存储服务启动
func StartNode(storageDir string, httpAddr string, clusterId string, clusterMembers string) error {
	if storageDir == "" {
		storageDir = DefaultStorageDir
	}

	if storageDir[len(storageDir)-1] != '/' {
		storageDir += "/"
	}

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

	node.Worker = worker.NewWorker(node.WorkerAddr, node.ClusterId)
	node.Worker.AddMates(strings.Split(node.ClusterMembers, ","))

	var err error
	node.NodeEnt, err = dao.GetNodeByWorkerId(node.Worker.Id) //&dao.Node{}
	if err != nil {
		return err
	}

	if node.NodeEnt == nil {
		node.NodeEnt = &dao.Node{
			NodeName: node.Worker.Id,
		}
		_, err = base.Insert(node.NodeEnt)
		if err != nil {
			return err
		}
	}

	node.NodeEnt.NodeAddr = node.WorkerAddr
	node.NodeEnt.Status = 0
	node.NodeEnt.TotalSpace = node.GetFreeSpace()
	node.NodeEnt.LastRegistered = time.Now()

	// 启动信息写入数据库
	err = node.NodeEnt.UpdateById([]string{"nodeAddr", "status", "totalSpace"})
	if err != nil {
		return err
	}

	node.Worker.OuterHandler = node

	// 准备启动服务
	serveErrorChan := make(chan error, 1)

	// start pingable server
	go func() {
		log.Debugf("http server will start at %v", node.WorkerAddr)
		serveErrorChan <- node.Worker.ServePingable()
	}()

	// start cluster node
	go func() {
		log.Debugf("worker server will start ")
		serveErrorChan <- node.Worker.Start()
	}()

	err = <-serveErrorChan
	log.Errorf("an error occur when serving storage: %v", err)

	return err
}

// 获取节点可用空闲空间
func (n *NodeStorage) GetFreeSpace() int64 {
	return 1024 * 1024 * 1024
}

// 将文件保存到本地 fileHash用于校验 未完成 暂无调用
// func (n *NodeStorage) SaveFile(filepath string, fileHash string) (*dao.FileIndex, error) {
// 	fileHandle, err := os.Open(filepath)
// 	if err != nil {
// 		log.Errorf("open %s failed:%v", filepath, err)
// 		return nil, err
// 	}

// 	defer fileHandle.Close()

// 	fileHash, err = filehash.CalcFileSha1(fileHandle)
// 	if err != nil {
// 		log.Errorf("calc filehash if uploaded file failed:%v", err)
// 		//StdErrResponse(c, ErrInternal)
// 		return nil, err
// 	}

// 	return &dao.FileIndex{
// 		FileName: filepath,
// 		FileHash: fileHash,
// 	}, nil

// }

// 从第一node复制过来备份
func (n *NodeStorage) SaveFileFromFileIndex(fileIndexId int, asNodeLevel string) (*dao.FileIndex, error) {

	//
	fileIndexEnt, err := dao.GetFileIndexById(fileIndexId)
	if err != nil {
		return nil, err
	}

	if fileIndexEnt == nil {
		return nil, nil
	}

	// 如果本地已经备份 则不用备份

	fileInStorage := n.FileHashToStoragePath(fileIndexEnt.FileHash)

	if DirExist(fileInStorage) {
		// 将已经备份的状态回写在文件记录上
		needUpCols := []string{}
		switch asNodeLevel {
		case "2":
			if fileIndexEnt.NodeId2 == 0 {
				fileIndexEnt.NodeId2 = n.NodeEnt.Id
				needUpCols = append(needUpCols, "nodeId2")
			}

		case "3":
			if fileIndexEnt.NodeId3 == 0 {
				fileIndexEnt.NodeId3 = n.NodeEnt.Id
				needUpCols = append(needUpCols, "nodeId3")
			}
		}
		if len(needUpCols) > 0 {
			err = fileIndexEnt.UpdateById(needUpCols)
			if err != nil {
				log.Errorf("when update fileIndex %d error:%v", fileIndexId, err)
			}
		}
		return fileIndexEnt, nil
	}

	// 获得外部访问别人的地址
	fileOutPath := n.FileHashToOutPath(fileIndexEnt.FileHash)
	// 外部访问带域名的地址
	fileServeUrl := n.WorkerAddr + fileOutPath
	log.Debugf("will save file form node1: %s", fileServeUrl)
	// 从对方节点下载到本地节点
	_, _, err = downloadFile(fileServeUrl, fileInStorage)
	if err != nil {
		log.Errorf("when download(%d) %s error:%v", fileIndexId, fileServeUrl, err)
		return nil, err
	}

	// 对比下载后的文件的hash
	realFileHash, _ := filehash.CalcSha1(fileInStorage)
	if realFileHash != fileIndexEnt.FileHash {
		// 当前计算的hash与数据库记录的不一致
		log.Errorf("fileIndex %d has wrong fileHash, realFileHash:%s, record:%s", fileIndexId, realFileHash, fileIndexEnt.FileHash)
		return nil, fmt.Errorf("fileIndex %d has wrong fileHash, realFileHash:%s, record:%s", fileIndexId, realFileHash, fileIndexEnt.FileHash)
	}

	// 按照备份级别，更新文件记录的备份状态
	needUpCols := []string{}
	switch asNodeLevel {
	case "2":
		if fileIndexEnt.NodeId2 == 0 {
			fileIndexEnt.NodeId2 = n.NodeEnt.Id
			needUpCols = append(needUpCols, "nodeId2")
		}

	case "3":
		if fileIndexEnt.NodeId3 == 0 {
			fileIndexEnt.NodeId3 = n.NodeEnt.Id
			needUpCols = append(needUpCols, "nodeId3")
		}
	}

	if len(needUpCols) > 0 {
		err = fileIndexEnt.UpdateById(needUpCols)
		if err != nil {
			log.Errorf("when update fileIndex %d error:%v", fileIndexId, err)
		}
	}

	return fileIndexEnt, nil
}

// 保存已经打开的文件流 必须在调用前知道fileHash // fileName 可空
func (n *NodeStorage) SaveFileHandler(inputFileHandler io.Reader, fileHash string, fileName string, size int64) (*dao.FileIndex, error) {

	// 新建本地规范路径的文件
	fileInStorage := n.FileHashToStoragePath(fileHash) //n.StorageDir + fileHash

	// 将handler层inputFileHandler 的保存在节点存储系统中
	fileInStorageHandle, err := os.OpenFile(fileInStorage, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		// 不能创建该文件，报错
		log.Errorf("open %s failed:%v", fileInStorage, err)
		return nil, err
	}

	defer fileInStorageHandle.Close()

	_, err = io.Copy(fileInStorageHandle, inputFileHandler)
	if err != nil {
		log.Errorf("file copy from uploaded tmp handle to %s failed:%v", fileInStorage, err)
		return nil, err
	}

	// 记录到数据库
	return n.collectFileInStorageToFileIndex(fileInStorage, fileHash, fileName, size)
}

// 将规范路径的本地文件记录到文件记录 fileHash必须提前计算好 fileName可空
func (n *NodeStorage) collectFileInStorageToFileIndex(filePathInStorage string, fileHash string, fileName string, fileSize int64) (*dao.FileIndex, error) {

	fileIndex := &dao.FileIndex{
		FileName:  fileName,
		FileHash:  fileHash,
		NodeId:    node.NodeEnt.Id,
		InnerPath: filePathInStorage,
		OuterPath: "/file/" + fileHash + "/" + fileName,
		Size:      fileSize,
		Space:     fileSize/1024 + 1,
		Status:    base.StatusNormal,
	}

	_, err := base.Insert(fileIndex)
	if err != nil {
		log.Errorf("insert %s failed:%v", filePathInStorage, err)
		return nil, err
	}

	log.Infof("a file stored, %s", filePathInStorage)

	// 分散告知其他节点
	// 获取按负载最低排序的节点
	condidateNodes := n.GetLowestRankedNodes()

	if len(condidateNodes) == 0 {
		log.Errorf("no mate members found, need to distribute other nodes")
	}

	// 选取2个负载最低的节点
	for i := 0; i < 2; i++ {
		mateId := condidateNodes[i].NodeName
		// 通知同伴节点备份文件
		n.DemandMateSaveFile(mateId, fileIndex.Id, strconv.Itoa(i+1))
	}

	return fileIndex, nil
}

func GetCurrentNode() *NodeStorage {
	return node
}

func (n *NodeStorage) FileHashToStoragePath(fileHash string) string {

	splitDeep := DirSplitDeep // 2
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

// 外网访问路径
func (n *NodeStorage) FileHashToOutPath(fileHash string) string {
	return "/file/" + fileHash
}

func (n *NodeStorage) ExecOfflineTask(task *dao.OfflineTask) (err error) {

	defer func() {
		// 延迟处理，将最后任务的状态保存到数据库
		if err != nil {
			task.Remark = err.Error()
			task.Status = dao.OfflineTaskStatusFail
		}
		task.UpdateById([]string{"status", "remark", "fileName", "size", "fileHash", "userFileId"})
	}()

	// 参数判断
	if task.Dataurl == "" {
		log.Errorf("dataurl is invalid when exec offlinetask")
		task.Status = dao.OfflineTaskStatusFail
		err = fmt.Errorf("dataurl is invalid ")
		return err
	}

	os.MkdirAll("./tmp/", os.ModeDir)

	localPath := fmt.Sprintf("./tmp/offline-%d", task.Id)

	// 下载任务指定的文件
	task.FileName, task.Size, err = downloadFile(task.Dataurl, localPath)
	if err != nil {
		log.Errorf("download dataurl(%s) error when exec offlinetask(%d): %v", task.Dataurl, task.Id, err)
		return err
	}

	log.Debugf("offline task %d has done. fileName:%s size:%d", task.Id, task.FileName, task.Size)

	if !DirExist(localPath) {
		err = fmt.Errorf("download offlinetask(%d) dataurl(%s) failed no local file", task.Id, task.Dataurl)
		log.Errorf(err.Error())
		return err
	}

	// 计算文件hash
	task.FileHash, err = filehash.CalcSha1(localPath)
	if err != nil {
		log.Errorf("calc filehash error when exec offlinetask(%d): %v", task.Id, err)
		return err
	}

	// 得到到本地规范目录路径
	fileInStorage := n.FileHashToStoragePath(task.FileHash)

	var fileIndex *dao.FileIndex
	// 保存到本地规范目录，并记录文件到数据库
	fileIndex, err = n.collectFileInStorageToFileIndex(fileInStorage, task.FileHash, task.FileName, task.Size)
	if err != nil {
		log.Errorf("offlinetask(%s) save to fileIndex error:%v", task.Id, err)
		return err
	}

	if fileIndex == nil {
		log.Errorf("offlinetask(%d) save to fileIndex failed, why?", task.Id)
		err = fmt.Errorf("offlinetask(%d) save to fileIndex failed, why?", task.Id)
		return err
	}

	// 获取保存文件的父目录
	parrentDir := "/"
	if task.ParentUserFileId > 0 {
		parrentUserFile, err := dao.GetUserFileById(task.ParentUserFileId)
		if err != nil {
			log.Warnf("get user(%d) parrentUserFile(%d) error:%v", task.UserId, task.ParentUserFileId)
		}
		if parrentUserFile != nil {
			parrentDir = parrentUserFile.FilePath + parrentUserFile.FileName
		}
	}

	// 生成用户文件记录
	userFile := &dao.UserFile{
		UserId:      task.UserId,
		FileIndexId: fileIndex.Id,
		FilePath:    parrentDir,
		FileHash:    task.FileHash,
		FileName:    task.FileName,
		Size:        task.Size,
		Space:       task.Size/1024 + 1,
	}

	// 生成用户文件的父目录hash 并保存
	userFile.MakePathHash()
	_, err = base.Insert(userFile)
	if err != nil {
		log.Errorf("save offlinetask(%d) to userFile error:%v", task.Id)
		return err
	}

	task.UserFileId = userFile.Id
	log.Debugf("offlinetask(%d) is done", task.Id)
	// 成功保存
	task.Status = dao.OfflineTaskStatusSaved

	return nil
}

// 获取按负载排序的节点列表
func (n *NodeStorage) GetLowestRankedNodes() []*dao.Node {
	nodeList := make([]*dao.Node, 0)
	nodeCondition := map[string]interface{}{
		"status=?": base.StatusNormal,
	}

	err := base.ListByCondition(&dao.Node{}, nodeCondition, 1, 1000, "unusedSpace desc", &nodeList)
	if err != nil {
		log.Errorf("list nodes failed:%v", err)
		return nil
	}

	// 转换成node的可排序对象
	var nodeListIf NodeList = nodeList
	sort.Sort(nodeListIf)

	return nodeList
}

/// 用于排序
type NodeList []*dao.Node

func (nl NodeList) Len() int {
	return len(nl)
}
func (nl NodeList) Less(i, j int) bool {
	return nl[i].UnusedSpace < nl[j].UnusedSpace
}
func (nl NodeList) Swap(i, j int) {
	nl[j], nl[i] = nl[i], nl[j]
}

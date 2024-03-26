package apihandler

import (
	"fmt"
	// "io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

type FileItem struct {
	Dirpath string `json:"dirpath"`
	Name    string `json:"name"`
	Url     string `json:"url"` //for view tpl mode, beego format this, vue us another format
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	Mtime   string `json:"mtime"`
	Thumb   string `json:"thumb"`
	// ThumbLoaded bool   `json:"-"`
}

var localDirRoot = "./"

var dirCache map[string][]*FileItem

func init() {
	dirCache = make(map[string][]*FileItem, 0)

	// 定时清理缓存
	go func() {
		tick := time.Tick(time.Second * 300)
		for _ = range tick {
			ClearCache()
		}
	}()
}

/*
*
dirpath must under localRoot
*/
func GetFileListFromDir(dirpath, dirPreRoute, filePreRoute string) []*FileItem {
	// dirCache = localDirPath=>[]*FileItem
	//log.Info("in GetFileItemListFromDir, dirpath=%s", dirpath)

	if dirpath == "" {
		//return nil
	}

	//curDirName := path.Base(dirpath)
	//parentDirName := path.Dir(dirpath)

	// set last char '/'
	if len(dirpath) > 0 && dirpath[len(dirpath)-1] != '/' {
		dirpath = dirpath + "/"
	}

	// return if exist
	if existFileItemList, ok := dirCache[dirpath]; ok {
		if existFileItemList != nil {
			return existFileItemList
		}
	}

	//log.Info("not exist:%s", dirpath)

	dirHandle, err := os.ReadDir(localDirRoot + "/" + dirpath)
	if err != nil {
		log.Warnf("open dir %s error:%v", dirpath, err)
		return nil
	}

	theDirList := make(FileItemSlice, 0)
	picIdx := 0
	//allNum := len(dirHandle)

	for _, fi := range dirHandle {

		lName := strings.ToLower(fi.Name())
		if len(fi.Name()) > 0 && fi.Name()[0] == '.' {
			continue
		}

		fileInfo, err := fi.Info()
		if err != nil {
			log.Warnf("get fileItem %s info error:%v, ignore", dirpath+fi.Name(), err)
			continue
		}

		if fi.IsDir() {
			// if lName == "thumbs" || lName == "thumb" {
			// 	continue
			// }

			// 目录 该目录下如果有封面，选出封面
			//thumbPath := GetThumbOfDir(dirpath+fi.Name(), filePreRoute)
			//log.Info("fi.name=%v thumb path=%v", fi.Name(), thumbPath)

			dirTitle := fi.Name() //getTitleOfDir(dirpath+fi.Name(), fi.Name()),//+fi.Name()+fmt.Sprintf("(%d/%d)", i+1, allNum)

			picItem := &FileItem{
				Dirpath: dirpath + fi.Name(),
				Name:    dirTitle,
				// Url:     dirPreRoute + "/" + dirpath + fi.Name(),
				IsDir: true,
				Size:  fileInfo.Size(),
				Mtime: fileInfo.ModTime().Format("2006-01-02 15:04"),
			}

			// if thumbPath := GetThumbOfDir(dirpath+fi.Name(), filePreRoute); thumbPath != "" {
			// 	picItem.Thumb = thumbPath
			// 	picItem.ThumbLoaded = true
			// }
			// go picItem.LoadThumb(filePreRoute)

			theDirList = append(theDirList, picItem)

		} else {
			// if lName == "thumb.jpg" || lName == "thumb.png" || lName == "thumb.gif" {
			// 	continue
			// }

			picIdx++

			// 只有图片才展示
			fExt := path.Ext(lName)
			thumbPath := dirpath + fi.Name()
			picItem := &FileItem{
				Dirpath: dirpath + fi.Name(),
				//Name:    fmt.Sprintf("%s-%d", curDirName, picIdx), //fmt.Sprintf("%s-%d", getTitleOfDir(dirpath, curDirName), picIdx),//
				Name:  fi.Name(),
				Url:   filePreRoute + "/" + dirpath + fi.Name(),
				Size:  fileInfo.Size(),
				Mtime: fileInfo.ModTime().Format("2006-01-02 15:04"),
			}
			if fExt == ".jpg" || fExt == ".png" || fExt == ".gif" {
				picItem.Thumb = filePreRoute + "/" + thumbPath
				// picItem.Url = filePreRoute + "/" + thumbPath
			}
			theDirList = append(theDirList, picItem)
		}
	}

	sort.Sort(theDirList)

	dirCache[dirpath] = theDirList
	log.Debugf("path %s is loaded into cache", dirpath)

	return theDirList
}

func ClearCache() {
	dirCache = make(map[string][]*FileItem, 0)
}

type FileItemSlice []*FileItem

func (ps FileItemSlice) Len() int {
	return len(ps)
}

func (ps FileItemSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// 对子目录排序
func (ps FileItemSlice) Less(i, j int) bool {

	iName, jName := path.Base(ps[i].Dirpath), path.Base(ps[j].Dirpath)
	if len(iName) > 0 && ('0' <= iName[0] && iName[0] <= '9') &&
		len(jName) > 0 && ('0' <= jName[0] && jName[0] <= '9') {

		in, jn := 0, 0
		fmt.Sscanf(iName, "%d", &in)
		fmt.Sscanf(jName, "%d", &jn)

		return in < jn
	}

	return iName < jName
}

type LocalFileListRequest struct {
	Dirpath  string `form:"dirpath"`
	FsRoute  string `form:"fsRoute"`
	Name     string `form:"fileName"`
	Page     int    `form:"page"`
	PageSize int    `form:"pagesize"`
}

const (
	FilePreRoute = "/file"
	pageSize     = 10
)

// FileIndexList - to list the FileIndex
func LocalFileList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &LocalFileListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		JsonErr(c, ErrInvalidParam)
		return
	}

	dirpath := c.Param("dirpath")

	if req.PageSize <= 0 {
		req.PageSize = pageSize
	}

	// 列表查询
	list := GetFileListFromDir(dirpath, "", FilePreRoute)

	if list == nil {
		log.Trace(requestId).Errorf("list failed, none")
		JsonErr(c, ErrInternal)
		return
	}

	allNum := len(list)

	p := paginator.NewPaginator(c.Request, req.PageSize, int64(allNum))
	// this.Data["paginator"] = p

	last := p.Page() * req.PageSize
	if last >= len(list) {
		last = len(list)
	}
	thePagedDirList := list[(p.Page()-1)*req.PageSize : last]

	JsonOk(c, gin.H{
		"total":     allNum,
		"page":      req.Page,
		"pagesize":  req.PageSize,
		"list":      thePagedDirList,
		"totalPage": p.PageNums(),
		"reqParam":  req,
		// "paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

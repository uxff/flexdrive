package handler

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func init() {
}

type UserFileListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	FileName    string `form:"fileName"`
	Dir         string `form:"dir"`
	SearchDir   int    `form:"searchDir"` // 是否在当前目录下搜索 默认搜索全部目录
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *UserFileListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.CreateStart != "" {
		condition["created>=?"] = r.CreateStart
	}

	if r.CreateEnd != "" {
		condition["created<=?"] = r.CreateEnd
	}

	if r.FileName != "" {
		// 关键词搜索模式
		condition["fileName like ?"] = "%" + r.FileName + "%"
		if r.SearchDir == 1 {
			fileIdxTmp := &dao.UserFile{
				FilePath: r.Dir,
			}
			condition["pathHash= ?"] = fileIdxTmp.MakePathHash()
			condition["filePath= ?"] = r.Dir
		}

	} else {
		// 目录浏览模式
		fileIdxTmp := &dao.UserFile{
			FilePath: r.Dir,
		}
		condition["pathHash= ?"] = fileIdxTmp.MakePathHash()
		condition["filePath= ?"] = r.Dir
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type UserFileItem struct {
	dao.UserFile
}

func NewUserFileItemFromEnt(fileIndexEnt *dao.UserFile) *UserFileItem {
	return &UserFileItem{
		UserFile: *fileIndexEnt,
	}
}

func UserFileList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &UserFileListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 当前浏览的目录
	if req.Dir == "" {
		req.Dir = "/"
	}
	// 补充尾巴的斜杠/
	if len(req.Dir) > 1 && req.Dir[len(req.Dir)-1] != '/' {
		req.Dir += "/"
	}

	condition := req.ToCondition()
	condition["status=?"] = base.StatusNormal // 只查询未删除

	// 列表查询
	list := make([]*dao.UserFile, 0)
	total, err := base.ListAndCountByCondition(&dao.UserFile{}, condition, req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*UserFileItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewUserFileItemFromEnt(v))
	}

	parentPath := path.Dir(strings.TrimRight(req.Dir, "/"))
	if len(parentPath) > 1 {
		parentPath += "/"
	}

	c.HTML(http.StatusOK, "userfile/list.tpl", gin.H{
		"LoginInfo":  getLoginInfo(c),
		"IsLogin":    isLoginIn(c),
		"total":      total,
		"page":       req.Page,
		"pagesize":   req.PageSize,
		"list":       resItems,
		"reqParam":   req,
		"dirLis":     NewDirLis(req.Dir),
		"parentPath": parentPath,
		"paginator":  paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

type DirLi struct {
	Dir    string
	Parent string
}

func NewDirLis(filePath string) []*DirLi {
	dirSlices := strings.Split(filePath, "/")
	dirLis := make([]*DirLi, 0)
	curParent := "/"
	for _, dirSlice := range dirSlices {
		if dirSlice == "" {
			continue
		}
		dirLis = append(dirLis, &DirLi{
			Dir:    dirSlice,
			Parent: curParent,
		})
		curParent += dirSlice + "/"
	}
	return dirLis
}

// 表示在ParentDir(必须已存在)下创建DirName(必须不存在)
type UserFileNewFolderRequest struct {
	ParentDir string `form:"parentDir" binding:"required"` // 左右都带/
	DirName   string `form:"dirName" binding:"required"`   // 不带/
}

func UserFileNewFolder(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &UserFileNewFolderRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	log.Trace(requestId).Debugf("req:%+v", req)

	userInfo := getLoginInfo(c)

	req.DirName = strings.Trim(req.DirName, "/ \t\n\r")

	if req.DirName == "" {
		StdErrMsgResponse(c, ErrInvalidParam, "要创建的文件夹名称为空，请重新提交")
		return
	}

	if req.ParentDir == "" {
		req.ParentDir = "/"
		// 根目录如果不存在，则创建根目录
		// parentDirEnt, _ := dao.GetUserFileByPath(userInfo.UserId, req.ParentDir)
		// if parentDirEnt != nil {
		// 	dao.MakeUserFilePath(userInfo.UserId, req.ParentDir, "/")
		// }

	} else {
		// 尾巴上补充上/
		if len(req.ParentDir) > 1 && req.ParentDir[len(req.ParentDir)-1] != '/' {
			req.ParentDir += "/"
		}
	}

	if strings.Contains(req.DirName, "/") {
		StdErrMsgResponse(c, ErrInvalidParam, "要创建的目录名称中不能有/等特殊字符")
		return
	}

	log.Trace(requestId).Debugf("will create path:%s  %s", req.ParentDir, req.DirName)

	// 父目录是否存在
	parentDirEnt, _ := dao.GetUserFileByPath(userInfo.UserId, req.ParentDir)
	if parentDirEnt == nil {
		StdErrMsgResponse(c, ErrInvalidParam, "选择的父目录("+req.ParentDir+")不存在")
		return
		//dao.MakeUserFilePath(userInfo.UserId, req.ParentDir)
	}

	// 当前目录如果存在则不创建

	existEnt, _ := dao.GetUserFileByPath(userInfo.UserId, req.ParentDir+req.DirName)
	if existEnt != nil {
		StdErrMsgResponse(c, ErrInternal, "选择的目录("+req.ParentDir+req.DirName+")已经存在")
		return
		//dao.MakeUserFilePath(userInfo.UserId, req.ParentDir)
	}

	// 创建该目录
	dirEnt, err := dao.MakeUserFilePath(userInfo.UserId, req.ParentDir, req.DirName)
	if err != nil {
		StdErrMsgResponse(c, ErrInvalidParam, "创建目录("+req.ParentDir+")失败")
		return
	}

	log.Trace(requestId).Debugf("创建目录%s成功", dirEnt.FilePath)

	c.Redirect(http.StatusMovedPermanently, RouteUserFileList+"?dir="+req.ParentDir)

}

func UserFileEnable(c *gin.Context) {
	fileIndexId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if fileIndexId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	//loginInfo := getLoginInfo(c)

	shareEnt, err := dao.GetUserFileById(int(fileIndexId))

	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if shareEnt == nil {
		StdErrResponse(c, ErrUserNotExist)
		return
	}

	if enable == 1 {
		// 启用
		shareEnt.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		shareEnt.Status = base.StatusDeleted
	}

	_, err = base.UpdateByCol("id", fileIndexId, shareEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteUserFileList)
}

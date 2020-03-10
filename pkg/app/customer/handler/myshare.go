package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func init() {
}

type ShareListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"name"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *ShareListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.CreateStart != "" {
		condition["created>=?"] = r.CreateStart
	}

	if r.CreateEnd != "" {
		condition["created<=?"] = r.CreateEnd
	}

	if r.Name != "" {
		condition["fileName like ?"] = "%" + r.Name + "%"
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type ShareItem struct {
	dao.Share
}

func NewShareItemFromEnt(shareEnt *dao.Share) *ShareItem {
	return &ShareItem{
		Share: *shareEnt,
	}
}

func ShareList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ShareListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Share, 0)
	total, err := base.ListAndCountByCondition(&dao.Share{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*ShareItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewShareItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "share/list.tpl", gin.H{
		"LoginInfo": getLoginInfo(c),
		"IsLogin":   isLoginIn(c),
		"total":     total,
		"page":      req.Page,
		"pagesize":  req.PageSize,
		"list":      resItems,
		"reqParam":  req,
		"paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

type ShareAddRequest struct {
	UserFileId  int       `form:"userFileId"`  // 文件
	ExpiredText string    `form:"expiredText"` // 有效期秒数
	ExpiredTime time.Time `form:"-"`
}

func ShareAdd(c *gin.Context) {
	// 已经分享的不做分享
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ShareAddRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		log.Trace(requestId).Errorf("bind param error:%v", err)
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	if req.UserFileId <= 0 {
		StdErrMsgResponse(c, ErrInvalidParam, "文件id为空")
		return
	}

	existShare, err := dao.GetShareByUserFile(req.UserFileId)
	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrMsgResponse(c, ErrInternal, "获取分享文件失败")
		return
	}

	loginInfo := getLoginInfo(c)

	userFile, err := dao.GetUserFileById(req.UserFileId)
	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrMsgResponse(c, ErrInternal, "查询分享文件失败")
		return
	}

	if userFile == nil {
		StdErrMsgResponse(c, ErrInternal, "要分享的文件不存在")
		return
	}

	if userFile.UserId != loginInfo.UserId {
		StdErrMsgResponse(c, ErrInternal, "只能分享自己的文件")
		return
	}

	if existShare != nil {
		//if existShare.Expired = // todo redirect
		StdResponse(c, ErrSuccess, existShare)
		return
	} else {
		shareItem := &dao.Share{
			UserId:     loginInfo.UserId,
			UserFileId: req.UserFileId,
			FileName:   userFile.FileName,
			FileHash:   userFile.FileHash,
			NodeId:     userFile.NodeId,
			//FilePath:   userFile.FilePath,
			Status: base.StatusNormal,
		}

		req.ExpiredTime, err = time.Parse("2006-01-02 15:04", req.ExpiredText)

		if req.ExpiredTime.Unix() > 0 {
			shareItem.Expired = req.ExpiredTime //time.Now().Add(time.Second * time.Duration(req.ExpiredSec))
		}

		_, err = base.Insert(shareItem)
		if err != nil {
			log.Trace(requestId).Errorf("db insert error:%v", err)
			StdErrMsgResponse(c, ErrInternal, "创建分享失败")
			return
		}
		StdResponse(c, ErrSuccess, shareItem) // todo redirect
		return
	}

}

// 查看用户文件是否已经被分享
func ShareCheck(c *gin.Context) {
	// 已经分享的不做分享
	requestId := c.GetString(CtxKeyRequestId)

	loginInfo := getLoginInfo(c)

	// 请求参数校验
	userFileId, _ := strconv.ParseInt(c.Param("userFileId"), 10, 64)
	if userFileId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	userFile, err := dao.GetUserFileById(int(userFileId))
	if err != nil || userFile == nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrMsgResponse(c, ErrInternal, "文件不存在或被删除")
		return
	}

	if userFile.UserId != loginInfo.UserId {
		StdErrMsgResponse(c, ErrInternal, "只能查看和分享自己的文件")
		return
	}

	existShare, err := dao.GetShareByUserFile(int(userFileId))
	if err != nil {
		log.Trace(requestId).Errorf("db error:%v", err)
		StdErrMsgResponse(c, ErrInternal, "查询分享失败")
		return
	}

	StdResponse(c, ErrSuccess, existShare)
	return
}

func ShareEnable(c *gin.Context) {
	shareId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if shareId <= 0 {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	shareEnt, err := dao.GetShareById(int(shareId))

	//_, err := base.GetByCol("id", mid, shareEnt)
	// exist, err := base.GetByCol("mid", mid, shareEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	if shareEnt == nil {
		StdErrResponse(c, ErrItemNotExist)
		return
	}

	if enable == 1 {
		// 启用
		shareEnt.Status = base.StatusNormal
	} else if enable == 9 {
		// 停用
		shareEnt.Status = base.StatusDeleted
	}

	//base.CacheDelByEntity("mgrLoginName", shareEnt.Email, shareEnt)

	_, err = base.UpdateByCol("id", shareId, shareEnt, []string{"status"})
	if err != nil {
		log.Errorf("db error:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	//StdResponse(c, ErrSuccess, nil)
	c.Redirect(http.StatusMovedPermanently, RouteShareList)
}

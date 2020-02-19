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
	UserFileId int `form:"userFileId"` // 文件
	ExpiredSec int `form:"expiredSec"` // 有效期秒数
}

func ShareAdd(c *gin.Context) {
	// 已经分享的不做分享
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ShareAddRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
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
		StdErrMsgResponse(c, ErrInternal, "要分享的文件不存在")
		return
	}

	if userFile.UserId != loginInfo.UserId {
		StdErrMsgResponse(c, ErrInternal, "只能分享自己的文件")
		return
	}

	if existShare != nil {
		//if existShare.Expired =
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

		if req.ExpiredSec > 0 {
			shareItem.Expired = time.Now().Add(time.Second * time.Duration(req.ExpiredSec))
		}

		_, err = base.Insert(shareItem)
		if err != nil {
			log.Trace(requestId).Errorf("db insert error:%v", err)
			StdErrMsgResponse(c, ErrInternal, "创建分享失败")
			return
		}
	}

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

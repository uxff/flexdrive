package apihandler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

func init() {
}

type ShareListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"name"`
	UserId      int    `form:"userId"`
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

	if r.UserId != 0 {
		condition["userId=?"] = r.UserId
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

// ShareList - to list share items
func ShareList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &ShareListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		JsonErr(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Share, 0)
	total, err := base.ListAndCountByCondition(&dao.Share{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*ShareItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewShareItemFromEnt(v))
	}

	JsonOk(c, gin.H{
		"total":    total,
		"page":     req.Page,
		"pagesize": req.PageSize,
		"list":     resItems,
		"reqParam": req,
		// "paginator": paginator.NewPaginator(c.Request, 10, int64(total)),
	})
}

// ShareEnable - to enable a share item
func ShareEnable(c *gin.Context) {
	shareId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if shareId <= 0 {
		JsonErr(c, ErrInvalidParam)
		return
	}

	enable, _ := strconv.ParseInt(c.Param("enable"), 10, 64)

	shareEnt, err := dao.GetShareById(int(shareId))

	//_, err := base.GetByCol("id", mid, shareEnt)
	// exist, err := base.GetByCol("mid", mid, shareEnt)
	if err != nil {
		log.Errorf("db error:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	if shareEnt == nil {
		JsonErr(c, ErrMgrNotExist)
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
		JsonErr(c, ErrInternal)
		return
	}

	JsonOk(c, gin.H{
		"id": shareEnt.Id,
	})
	// c.Redirect(http.StatusMovedPermanently, RouteShareList)
}

// ShareDetail - to show detail of a share item
func ShareDetail(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	shareHash := c.Param("shareHash")
	if shareHash == "" {
		JsonErr(c, ErrInvalidParam)
		return
	}

	shareItem, err := dao.GetShareByShareHash(shareHash)
	if err != nil {
		log.Trace(requestId).Debugf("get shareHash(%s) error:%v", shareHash, err)
		JsonErr(c, ErrInternal)
		return
	}

	if shareItem == nil || shareItem.Status == base.StatusDeleted {
		JsonErrMsg(c, ErrItemNotExist, "分享的内容不存在或已删除")
		return
	}

	genShareOutPath(c, shareItem)

	JsonOk(c, shareItem)
}

func genShareOutPath(c *gin.Context, i *dao.Share) string {
	i.MakeShareHash()
	i.OuterPath = "http://" + c.Request.Host + "/s/" + i.ShareHash
	return i.OuterPath
}

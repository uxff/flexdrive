package apihandler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils"
)

func init() {
}

type NodeListRequest struct {
	Status     int    `json:"status"`
	Name       string `json:"name"`
	LastActive int    `json:"lastActive"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pagesize"`
}

func (r *NodeListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.Status > 0 {
		condition["status = ?"] = r.Status
	}

	if r.Name != "" {
		condition["nodeName like ?"] = "%" + r.Name + "%"
	}

	if r.LastActive > 0 {
		condition["lastRegistered > ?"] = time.Now().Add(-time.Duration(r.LastActive) * time.Second).Format("2006-01-02 15:04:05")
	}

	log.Debugf("r=%+v tocondition:%+v", r, condition)
	return condition
}

// 接口返回的元素
type NodeItem struct {
	dao.Node
	TotalSpaceDesc  string
	UsedSpaceDesc   string
	UnusedSpaceDesc string
}

func NewNodeItemFromEnt(nodeEnt *dao.Node) *NodeItem {
	return &NodeItem{
		Node:            *nodeEnt,
		TotalSpaceDesc:  utils.SizeToHuman(nodeEnt.TotalSpace),
		UsedSpaceDesc:   utils.SizeToHuman(nodeEnt.UsedSpace),
		UnusedSpaceDesc: utils.SizeToHuman(nodeEnt.UnusedSpace),
	}
}

func NodeList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &NodeListRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		JsonErr(c, ErrInvalidParam+":"+err.Error())
		return
	}

	// 列表查询
	list := make([]*dao.Node, 0)
	total, err := base.ListAndCountByCondition(&dao.Node{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		JsonErr(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*NodeItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewNodeItemFromEnt(v))
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

func NodeSetCapacity(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)
	nodeIdStr := c.PostForm("nodeId")
	totalSpaceStr := c.PostForm("totalSpace")
	nodeId, _ := strconv.Atoi(nodeIdStr)
	totalSpace, _ := strconv.Atoi(totalSpaceStr)
	if nodeId <= 0 {
		JsonErrMsg(c, ErrInternal, "没有提交节点id")
		return
	}
	nodeEnt, err := dao.GetNodeById(nodeId)
	if err != nil {
		log.Trace(requestId).Errorf("get node(%d) error:%v", nodeId, err)
		JsonErrMsg(c, ErrInternal, "节点id查询失败:"+err.Error())
		return
	}
	if nodeEnt == nil {
		JsonErrMsg(c, ErrInternal, "节点id不存在")
		return
	}
	if int64(totalSpace) < nodeEnt.UsedSpace {
		JsonErrMsg(c, ErrInternal, "空间不能小于节点的已用空间")
		return
	}

	nodeEnt.TotalSpace = int64(totalSpace)
	nodeEnt.UnusedSpace = nodeEnt.TotalSpace - nodeEnt.UsedSpace

	err = nodeEnt.UpdateById([]string{"totalSpace", "unusedSpace"})
	if err != nil {
		log.Trace(requestId).Errorf("update node(%d) error:%v", nodeId, err)
		JsonErrMsg(c, ErrInternal, "节点id更新错误:"+err.Error())
		return
	}

	JsonOk(c, gin.H{
		"id": nodeEnt.Id,
	})
	// c.Redirect(http.StatusMovedPermanently, RouteNodeList)
}

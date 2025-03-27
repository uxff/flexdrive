package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uxff/flexdrive/pkg/app/nodestorage/model/storagemodel"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
	"github.com/uxff/flexdrive/pkg/utils/paginator"
)

func init() {
}

type NodeListRequest struct {
	CreateStart string `form:"createStart"`
	CreateEnd   string `form:"createEnd"`
	Name        string `form:"name"`
	LastActive  int    `form:"lastActive"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pagesize"`
}

func (r *NodeListRequest) ToCondition() (condition map[string]interface{}) {
	condition = make(map[string]interface{})

	if r.CreateStart != "" {
		condition["created>=?"] = r.CreateStart
	}

	if r.CreateEnd != "" {
		condition["created<=?"] = r.CreateEnd
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
}

func NewNodeItemFromEnt(nodeEnt *dao.Node) *NodeItem {
	return &NodeItem{
		Node: *nodeEnt,
	}
}

func NodeList(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)

	// 请求参数校验
	req := &NodeListRequest{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		StdErrResponse(c, ErrInvalidParam)
		return
	}

	// 列表查询
	list := make([]*dao.Node, 0)
	total, err := base.ListAndCountByCondition(&dao.Node{}, req.ToCondition(), req.Page, req.PageSize, "", &list)
	if err != nil {
		log.Trace(requestId).Errorf("list failed:%v", err)
		StdErrResponse(c, ErrInternal)
		return
	}

	// 从数据库结构转换成返回结构
	resItems := make([]*NodeItem, 0)
	for _, v := range list {
		resItems = append(resItems, NewNodeItemFromEnt(v))
	}

	c.HTML(http.StatusOK, "node/list.tpl", gin.H{
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

func NodeSetspace(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)
	nodeIdStr := c.PostForm("nodeId")
	totalSpaceStr := c.PostForm("totalSpace")
	nodeId, _ := strconv.Atoi(nodeIdStr)
	totalSpace, _ := strconv.Atoi(totalSpaceStr)
	if nodeId <= 0 {
		StdErrMsgResponse(c, ErrInternal, "没有提交节点id")
		return
	}
	nodeEnt, err := dao.GetNodeById(nodeId)
	if err != nil {
		log.Trace(requestId).Errorf("get node(%d) error:%v", nodeId, err)
		StdErrMsgResponse(c, ErrInternal, "节点id查询失败:"+err.Error())
		return
	}
	if nodeEnt == nil {
		StdErrMsgResponse(c, ErrInternal, "节点id不存在")
		return
	}
	if int64(totalSpace) < nodeEnt.UsedSpace {
		StdErrMsgResponse(c, ErrInternal, "空间不能小于节点的已用空间")
		return
	}

	nodeEnt.TotalSpace = int64(totalSpace)
	nodeEnt.UnusedSpace = nodeEnt.TotalSpace - nodeEnt.UsedSpace

	err = nodeEnt.UpdateById([]string{"totalSpace", "unusedSpace"})
	if err != nil {
		log.Trace(requestId).Errorf("update node(%d) error:%v", nodeId, err)
		StdErrMsgResponse(c, ErrInternal, "节点id更新错误:"+err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, RouteNodeList)
}

// kick a node in the cluster
func NodeKick(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)
	nodeIdStr := c.Param("id") // id in the db, not clusterMemberId
	nodeId, _ := strconv.Atoi(nodeIdStr)

	if nodeId <= 0 {
		StdErrMsgResponse(c, ErrInternal, "没有提交节点id")
		return
	}

	nodeEnt, err := dao.GetNodeById(nodeId)
	if err != nil {
		log.Trace(requestId).Errorf("get node(%d) error:%v", nodeId, err)
		StdErrMsgResponse(c, ErrInternal, "节点id查询失败:"+err.Error())
		return
	}
	if nodeEnt == nil {
		StdErrMsgResponse(c, ErrInternal, "节点id不存在")
		return
	}

	err = storagemodel.KickMember(nodeEnt.NodeAddr)
	if err != nil {
		log.Trace(requestId).Errorf("kick node(%d) error:%v", nodeId, err)
	}

	c.Redirect(http.StatusMovedPermanently, RouteNodeList)
}

// invite node join the cluster
func NodeInvite(c *gin.Context) {
	requestId := c.GetString(CtxKeyRequestId)
	nodeIdStr := c.Param("id") // id in the db, not clusterMemberId
	nodeId, _ := strconv.Atoi(nodeIdStr)

	if nodeId <= 0 {
		StdErrMsgResponse(c, ErrInternal, "没有提交节点id")
		return
	}

	nodeEnt, err := dao.GetNodeById(nodeId)
	if err != nil {
		log.Trace(requestId).Errorf("get node(%d) error:%v", nodeId, err)
		StdErrMsgResponse(c, ErrInternal, "节点id查询失败:"+err.Error())
		return
	}
	if nodeEnt == nil {
		StdErrMsgResponse(c, ErrInternal, "节点id不存在")
		return
	}

	err = storagemodel.AddMember(nodeEnt.NodeAddr)
	if err != nil {
		log.Trace(requestId).Errorf("invite node(%d) error:%v", nodeId, err)
	}

	c.Redirect(http.StatusMovedPermanently, RouteNodeList)
}

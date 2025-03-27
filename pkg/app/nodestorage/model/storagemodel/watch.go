package storagemodel

import (
	"time"

	"github.com/uxff/flexdrive/pkg/app/nodestorage/clusterworker"
	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

const (
	RegisterTimeoutSec  = 300
	RegisterIntervalSec = 30
)

var presetNodeList []*clusterworker.Worker

// watch and clear lost mates in db
func (n *NodeStorage) WatchMates() error {
	return n.WatchMatesFronDB()
}
func (n *NodeStorage) WatchMatesFromEnv() error {
	return nil
}
func (n *NodeStorage) WatchMatesFronDB() error {
	for {
		// register self
		n.NodeEnt.LastRegistered = time.Now()
		n.NodeEnt.Status = base.StatusNormal
		n.NodeEnt.Follow = n.Worker.GetMasterId()
		n.NodeEnt.UpdateById([]string{"lastRegistered", "status", "follow"})

		// check mates are registered
		nodeList := make([]dao.Node, 0)
		condition := map[string]interface{}{
			"clusterId=?": n.ClusterId,
			"status=?":    base.StatusNormal,
		}
		err := base.ListByCondition(&dao.Node{}, condition, 1, 100000, "", &nodeList)
		if err != nil {
			log.Errorf("list nodes failed:%v", err)
			return err
		}

		for _, mate := range nodeList {
			if time.Now().Unix()-mate.LastRegistered.Unix() > 300 {
				mate.Status = base.StatusInactive
				mate.UpdateById([]string{"status"})
				log.Debugf("a mate is down: %s", mate.NodeAddr)
				//n.Worker.KickMate(mate.Id)//todo locked
				// delete(n.Worker.ClusterMembers, mate.NodeName)
				// n.Worker.ClusterMembersMap.Delete(mate.NodeName)
			} else {
				// log.Debugf("detected a mate(%s):%s ", mate.NodeName, mate.NodeAddr)
				// n.Worker.AddMates([]string{mate.NodeAddr}) // 不能从mysql数据库来决定集群节点数量
			}
		}

		time.Sleep(RegisterIntervalSec * time.Second)
	}
}

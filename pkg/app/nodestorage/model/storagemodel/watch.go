package storagemodel

import (
	"time"

	"github.com/uxff/flexdrive/pkg/dao"
	"github.com/uxff/flexdrive/pkg/dao/base"
	"github.com/uxff/flexdrive/pkg/log"
)

const (
	RegisterTimeoutSec = 300
)

// watch and clear lost mates in db
func (n *NodeStorage) WatchMates() error {
	for {
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
				mate.Status = base.StatusDeleted
				mate.UpdateById([]string{"status"})
				log.Debugf("a mate is down: %s", mate.NodeAddr)
				//n.Worker.KickMate(mate.Id)//todo locked
				delete(n.Worker.ClusterMembers, mate.NodeName)
			} else {
				log.Debugf("detected a mate(%s):%s ", mate.NodeName, mate.NodeAddr)
				n.Worker.AddMates([]string{mate.NodeAddr})
			}
		}

		time.Sleep(1 * time.Second)
	}
}

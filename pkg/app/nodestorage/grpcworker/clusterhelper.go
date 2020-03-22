package grpcworker

import (
	"crypto/md5"
	"encoding/hex"
)

const regPreKey = "/mycluster/"
const masterKey = "master"
const membersKey = "workers"

type clusterHelper struct {
	clusterId string
}

func NewClusterHelper(clusterId string) *clusterHelper {
	return &clusterHelper{
		clusterId,
	}
}

func (c *clusterHelper) genClusterKeyBase() string {
	return regPreKey + c.clusterId + "/"
}

//// 生成一个字符串key(比如redis key) 用于保存 workers
//func (c *clusterHelper) genMasterIdKey() string {
//	return c.genClusterKeyBase() + masterKey
//}

//// 生成一个字符串key(比如redis key) 用于保存 worders
//func (c *clusterHelper) genMembersKey() string {
//	return c.genClusterKeyBase() + membersKey
//}

func (c *clusterHelper) genMemberHash(memberId string) string {
	s := c.genClusterKeyBase() + memberId
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

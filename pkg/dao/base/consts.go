package base

// 行数据的状态
const (
	StatusNormal  = 1
	StatusDeleted = 99
)

const (
	// 小于10的mid是超级管理 不可删除
	MaxSuperMid = 10
	// 小于10的roleid是超级管理角色 不可删除
	MaxSuperRoleId = 10
)

var StatusMap = map[int]string{
	StatusNormal:  "正常",
	StatusDeleted: "已删除",
}

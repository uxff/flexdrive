package rbac

import (
	"strings"
)

type ResourceItem struct {
	Name string `json:"name"`
	//Route    string     `json:"route"` // 菜单展示时候的前端路由 权限不用
	PermitRoute string     `json:"permitRoute"` // 后端权限使用 表示前端页面下的调用接口使用路由
	Icon        string     `json:"icon"`
	Access      bool       `json:"access"`
	Sub         RoleAccess `json:"sub"`
}

// 一组资源 以json存在roles.RAccess字段中
type RoleAccess []*ResourceItem

// 从json转换为资源组
// func (r *RoleAccess) FromString(str string) {
// 	json.Unmarshal([]byte(str), r)
// }

// // 将资源组转换为json
// func (r *RoleAccess) ToString() string {
// 	b, _ := json.Marshal(r)
// 	return string(b)
// }

// 检查资源组中有没有路由权限
func (r *RoleAccess) CheckRouteAccessable(apiRouteStr string) bool {
	for _, rsItem := range *r {
		if rsItem.PermitRoute != "" && strings.HasPrefix(apiRouteStr+"/", rsItem.PermitRoute+"/") && rsItem.Access == true {
			return true
		}

		if rsItem.Sub.CheckRouteAccessable(apiRouteStr) {
			return true
		}
	}
	return false
}

// 获取所有菜单并注明是否有权限 获取角色的资源组
func GetAllAccessItems(customizedRoleAccess RoleAccess) RoleAccess {
	// 模板menu
	allItem := GetAllMenu()
	matchAccessWithCustomizedItems(allItem, customizedRoleAccess)
	return allItem
}

// 从保存定义过权限中，补充到全部菜单(通过getAllMune获得的模板菜单)上
func matchAccessWithCustomizedItems(allItem, customizedRoleAccess RoleAccess) {
	for _, aItem := range allItem {
		//
		aItem.Access = customizedRoleAccess.CheckRouteAccessable(aItem.PermitRoute)
		if aItem.Sub != nil {
			matchAccessWithCustomizedItems(aItem.Sub, customizedRoleAccess)
		}
	}
}

// 基于路由返回是否有权限
// func CheckAccessByRoute(customizedRoleAccess RoleAccess, apiRouteStr string) bool {
// 	allRsItem := GetAllAccessItems(customizedRoleAccess)
// 	return allRsItem.CheckRouteAccessable(apiRouteStr)
// }

// 10 以内的roleid是超级管理 拥有全部权限
func GetMaxSuperRoleId() int {
	return 10
}

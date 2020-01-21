package rbac

import (
	"encoding/json"
	//"gitlab.acewill.cn/rpc/GoPay/apps/admin/model/menu"
	"strings"

	"github.com/uxff/flexdrive/pkg/dao"
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

// 将菜单组复制为资源组
// func newResourceItemFromMenuItem(m *menu.MenuItem) *ResourceItem {
// 	r := &ResourceItem{}
// 	r.Name = m.Name
// 	r.Route = m.Route
// 	r.PermitRoute = m.PermitRoute
// 	r.Icon = m.Icon
// 	for _, mSub := range m.Sub {
// 		if mSub != nil {
// 			if rSub := newResourceItemFromMenuItem(mSub); rSub != nil {
// 				r.Sub = append(r.Sub, rSub)
// 			}
// 		}
// 	}
// 	return r
// }

// 将资源组复制为菜单组
// func newMenuItemFromResourceItem(r *ResourceItem, onlyAccessed bool) *menu.MenuItem {
// 	if onlyAccessed && !r.Access {
// 		return nil
// 	}
// 	m := &menu.MenuItem{}
// 	m.Name = r.Name
// 	m.Route = r.Route
// 	m.PermitRoute = r.PermitRoute
// 	m.Icon = r.Icon
// 	for _, rSub := range r.Sub {
// 		if rSub != nil {
// 			if mSub := newMenuItemFromResourceItem(rSub, onlyAccessed); mSub != nil {
// 				m.Sub = append(m.Sub, mSub)
// 			}
// 		}
// 	}
// 	return m
// }

// 从json转换为资源组
func (r *RoleAccess) FromString(str string) {
	json.Unmarshal([]byte(str), r)
}

// 将资源组转换为json
func (r *RoleAccess) ToString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// 检查资源组中有没有路由权限
func (r *RoleAccess) CheckRouteAccessable(apiRouteStr string) bool {
	for _, rsItem := range *r {
		if rsItem.PermitRoute != "" && strings.HasPrefix(apiRouteStr, rsItem.PermitRoute) && rsItem.Access == true {
			return true
		}

		if rsItem.Sub.CheckRouteAccessable(apiRouteStr) {
			return true
		}
	}
	return false
}

// 获取所有菜单并注明是否有权限 获取角色的资源组
func GetAccessMenuByRoleEnt(role *dao.Role) RoleAccess {
	allItem := GetAllMenu()
	//allItem := make(RoleAccess, 0) // 有权限的组

	// 全部菜单初始化 为空
	// for _, menu := range allMenu {
	// 	allItem = append(allItem, newResourceItemFromMenuItem(menu))
	// }

	// 数据库保存的某个角色的定义过的权限
	customizedRoleAccess := make(RoleAccess, 0)

	if role.RAccess != "" {
		customizedRoleAccess.FromString(role.RAccess)
	}

	matchAccessWithCustomizedItems(allItem, customizedRoleAccess)

	// return customizedRoleAccess
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

// // 基于角色权限获取菜单
// func GetMenuByRoleEnt(role *dao.Role) []*ResourceItem {
// 	allRsItem := GetAccessMenuByRoleEnt(role)
// 	//allMenuItem := make([]*menu.MenuItem, 0)

// 	for _, rsItem := range allRsItem {
// 		rsItem.Access = true // 一级菜单是有权限的
// 		mItem := newMenuItemFromResourceItem(rsItem, true)
// 		if mItem != nil && len(mItem.Sub) > 0 {
// 			allMenuItem = append(allMenuItem, mItem)
// 		}
// 	}

// 	return allMenuItem
// }

// 基于路由返回是否有权限
func CheckAccessByRoute(role *dao.Role, apiRouteStr string) bool {
	allRsItem := GetAccessMenuByRoleEnt(role)
	return allRsItem.CheckRouteAccessable(apiRouteStr)
}

// 10 以内的roleid是超级管理 拥有全部权限
func GetMaxSuperRoleId() int {
	return 10
}

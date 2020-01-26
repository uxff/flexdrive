package rbac

var defaultMenus = RoleAccess{
	&ResourceItem{
		Name:        "首页",
		Icon:        "",
		PermitRoute: "/",
	},
	&ResourceItem{
		Name:        "节点管理",
		Icon:        "node",
		PermitRoute: "/node/list", // 菜单组的虚拟api路由 没用 需要非空
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "节点列表",
				PermitRoute: "/api/merchant",
			},
		},
	},
	&ResourceItem{
		Name:        "文件管理",
		Icon:        "file",
		PermitRoute: "/api/set", // 菜单组的虚拟api路由 没用 非空就可以
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "文件列表",
				PermitRoute: "/api/role",
			},
		},
	},
	&ResourceItem{
		Name:        "会员管理",
		Icon:        "customer",
		PermitRoute: "/api/set", // 菜单组的虚拟api路由 没用 非空就可以
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "会员列表",
				PermitRoute: "/api/role",
			},
			&ResourceItem{
				Name:        "会员等级",
				PermitRoute: "/api/manager",
			},
			&ResourceItem{
				Name:        "订单管理",
				PermitRoute: "/api/manager",
			},
			&ResourceItem{
				Name:        "分享管理",
				PermitRoute: "/api/manager",
			},
		},
	},
	&ResourceItem{
		Name:        "系统管理",
		Icon:        "setting",
		PermitRoute: "/api/set", // 菜单组的虚拟api路由 没用 非空就可以
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "角色管理",
				PermitRoute: "/api/role",
			},
			&ResourceItem{
				Name:        "管理员账号管理",
				PermitRoute: "/api/manager",
			},
		},
	},
}

func GetAllMenu() RoleAccess {
	return defaultMenus
}

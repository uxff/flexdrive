package rbac

// 权限资源列表 不能当菜单
var defaultMenus = RoleAccess{
	&ResourceItem{
		Name:        "节点管理",
		Icon:        "glyphicon-th-large",
		PermitRoute: "/node/list", // 菜单组的虚拟api路由 没用 需要非空
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "节点列表",
				PermitRoute: "/node/list",
			},
		},
	},
	&ResourceItem{
		Name:        "文件管理",
		Icon:        "file",
		PermitRoute: "/file", // 菜单组的虚拟api路由 没用 非空就可以
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "文件列表",
				PermitRoute: "/file",
			},
		},
	},
	&ResourceItem{
		Name:        "会员管理",
		Icon:        "glyphicon-list-alt",
		PermitRoute: "/user", // 菜单组的虚拟api路由 没用 非空就可以
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "会员列表",
				PermitRoute: "/user",
			},
			&ResourceItem{
				Name:        "会员等级",
				PermitRoute: "/userlevel",
			},
			&ResourceItem{
				Name:        "订单管理",
				PermitRoute: "/order",
			},
			&ResourceItem{
				Name:        "分享管理",
				PermitRoute: "/share",
			},
		},
	},
	&ResourceItem{
		Name:        "系统管理",
		Icon:        "glyphicon-list-alt",
		PermitRoute: "/", // 菜单组的虚拟api路由 没用 非空就可以
		Sub: []*ResourceItem{
			&ResourceItem{
				Name:        "角色管理",
				PermitRoute: "/role/list",
			},
			&ResourceItem{
				Name:        "角色添加",
				PermitRoute: "/role/add",
			},
			&ResourceItem{
				Name:        "角色编辑及权限设置",
				PermitRoute: "/role/edit",
			},
			&ResourceItem{
				Name:        "角色启用停用",
				PermitRoute: "/role/enable",
			},
			&ResourceItem{
				Name:        "管理员账号管理",
				PermitRoute: "/manager/list",
			},
			&ResourceItem{
				Name:        "管理员账号添加",
				PermitRoute: "/manager/add",
			},
			&ResourceItem{
				Name:        "管理员账号编辑",
				PermitRoute: "/manager/edit",
			},
			&ResourceItem{
				Name:        "管理员账号启用停用",
				PermitRoute: "/manager/enable",
			},
		},
	},
}

func GetAllMenu() RoleAccess {
	return defaultMenus
}

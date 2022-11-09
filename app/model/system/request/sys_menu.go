package request

import (
	"time"
)

type SysMenu struct {
	OpenPage   bool       `json:"openPage"`                 // 开启分页
	PageNum    int        `json:"pageNum" form:"pageNum"`   // 页码
	PageSize   int        `json:"pageSize" form:"pageSize"` // 每页大小
	MenuID     int64      `json:"menuId"`                   // 菜单ID
	MenuName   string     `json:"menuName"`                 // 菜单名称
	ParentID   int64      `json:"parentId"`                 // 父菜单ID
	OrderNum   int64      `json:"orderNum"`                 // 显示顺序
	Path       string     `json:"path"`                     // 路由地址
	Component  string     `json:"component"`                // 组件路径
	Query      string     ` json:"query"`                   // 路由参数
	IsFrame    int64      `json:"isFrame"`                  // 是否为外链（0是 1否）
	IsCache    int64      ` json:"isCache"`                 // 是否缓存（0缓存 1不缓存）
	MenuType   string     `json:"menuType"`                 // 菜单类型（M目录 C菜单 F按钮）
	Visible    string     ` json:"visible"`                 // 菜单状态（0显示 1隐藏）
	Status     string     `json:"status"`                   // 菜单状态（0正常 1停用）
	Perms      string     `json:"perms"`                    // 权限标识
	Icon       string     `json:"icon"`                     // 菜单图标
	CreateBy   string     `json:"createBy"`                 // 创建者
	CreateTime time.Time  `json:"createTime"`               // 创建时间
	UpdateBy   string     `json:"updateBy"`                 // 更新者
	UpdateTime time.Time  `json:"updateTime"`               // 更新时间
	Remark     string     `json:"remark"`                   // 备注
	Children   []*SysMenu `json:"children,omitempty"`       // 子菜单
	ArrIdx     int        `json:"arrIdx,omitempty"`         // 临时变量 用于生成tree
	RoleId     int64      `json:"roleId"`                   // 角色Id
	UserId     int64      `json:"userId"`                   // 用户Id
}

type MenuListByRoleId struct {
	RoleId              int64 `json:"roleId"`
	IsMenuCheckStrictly bool  `json:"isMenuCheckStrictly"`
}

package response

import (
	"go-vea/app/model/system"
	"time"
)

type SysMenu struct {
	MenuID     int64     `json:"menuId"`     // 菜单ID
	MenuName   string    `json:"menuName"`   // 菜单名称
	ParentID   int64     `json:"parentId"`   // 父菜单ID
	OrderNum   int64     `json:"orderNum"`   // 显示顺序
	Path       string    `json:"path"`       // 路由地址
	Component  string    `json:"component"`  // 组件路径
	Query      string    `json:"query"`      // 路由参数
	IsFrame    int64     `json:"isFrame"`    // 是否为外链（0是 1否）
	IsCache    int64     `json:"isCache"`    // 是否缓存（0缓存 1不缓存）
	MenuType   string    `json:"menuType"`   // 菜单类型（M目录 C菜单 F按钮）
	Visible    string    `json:"visible"`    // 菜单状态（0显示 1隐藏）
	Status     string    `json:"status"`     // 菜单状态（0正常 1停用）
	Perms      string    `json:"perms"`      // 权限标识
	Icon       string    `json:"icon"`       // 菜单图标
	CreateBy   string    `json:"createBy"`   // 创建者
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateBy   string    `json:"updateBy"`   // 更新者
	UpdateTime time.Time `json:"updateTime"` // 更新时间
	Remark     string    `json:"remark"`     // 备注
	Children   []SysMenu `json:"children"`   // 子菜单
}

// RouterVo 路由配置信息
type RouterVo struct {
	Name       string     `json:"name"`                 // 路由名字
	Path       string     `json:"path"`                 // 路由地址
	Hidden     bool       `json:"hidden"`               // 是否隐藏路由，当设置 true 的时候该路由不会再侧边栏出现
	Redirect   string     `json:"redirect,omitempty"`   // 重定向地址，当设置 noRedirect 的时候该路由在面包屑导航中不可被点击
	Component  string     `json:"component"`            // 组件地址
	Query      string     `json:"query,omitempty"`      // 路由参数：如 {"id": 1, "name": "admin"}
	AlwaysShow bool       `json:"alwaysShow,omitempty"` // 当你一个路由下面的 children 声明的路由大于1个时，自动会变成嵌套的模式--如组件页面
	Meta       *MetaVo    `json:"meta"`                 // 其他元素
	Children   []RouterVo `json:"children,omitempty"`   // 子路由
}

// MetaVo 路由显示信息
type MetaVo struct {
	Title   string `json:"title"`   // 设置该路由在侧边栏和面包屑中展示的名字
	Icon    string `json:"icon"`    // 设置该路由的图标，对应路径src/assets/icons/svg
	NoCache bool   `json:"noCache"` // 设置为true，则不会被 <keep-alive>缓存
	Link    string `json:"link"`    // 内链地址（http(s)://开头）
}

type RoleMenuTreeSelectResp struct {
	CheckedKeys []int64           `json:"checkedKeys"`
	Menus       []*system.SysMenu `json:"menus"`
}

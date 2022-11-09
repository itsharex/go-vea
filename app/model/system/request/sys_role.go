package request

import (
	"time"
)

type SysRole struct {
	OpenPage          bool      `json:"openPage"`                 // 开启分页
	PageNum           int       `json:"pageNum" form:"pageNum"`   // 页码
	PageSize          int       `json:"pageSize" form:"pageSize"` // 每页大小
	Ids               []int64   `json:"ids"`                      // roleIds
	RoleID            int64     `json:"roleId"`                   // 角色ID
	RoleName          string    `json:"roleName"`                 // 角色名称
	RoleKey           string    `json:"roleKey"`                  // 角色权限字符串
	RoleSort          int64     `json:"roleSort"`                 // 显示顺序
	DataScope         string    `json:"dataScope"`                // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）
	MenuCheckStrictly bool      `json:"menuCheckStrictly"`        // 菜单树选择项是否关联显示
	DeptCheckStrictly bool      `json:"deptCheckStrictly"`        // 部门树选择项是否关联显示
	Status            string    `json:"status"`                   // 角色状态（0正常 1停用）
	DelFlag           string    `json:"delFlag"`                  // 删除标志（0代表存在 2代表删除）
	CreateBy          string    `json:"createBy"`                 // 创建者
	CreateTime        time.Time `json:"createTime"`               // 创建时间
	UpdateBy          string    `json:"updateBy"`                 // 更新者
	UpdateTime        time.Time `json:"updateTime"`               // 更新时间
	Remark            string    `json:"remark"`                   // 备注
}

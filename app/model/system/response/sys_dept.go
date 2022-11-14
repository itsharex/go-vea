package response

import (
	"go-vea/app/common/base"
	"go-vea/app/model/system"
	"time"
)

type SysDept struct {
	base.CommonModel
	DeptID     int64     `json:"deptId"`     // 部门id
	ParentID   int64     `json:"parentId"`   // 父部门id
	Ancestors  string    `json:"ancestors"`  // 祖级列表
	DeptName   string    `json:"deptName"`   // 部门名称
	OrderNum   int64     `json:"orderNum"`   // 显示顺序
	Leader     string    `json:"leader"`     // 负责人
	Phone      string    `json:"phone"`      // 联系电话
	Email      string    `json:"email"`      // 邮箱
	Status     string    `json:"status"`     // 部门状态（0正常 1停用）
	DelFlag    string    `json:"delFlag"`    // 删除标志（0代表存在 2代表删除）
	CreateBy   string    `json:"createBy"`   // 创建者
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateBy   string    `json:"updateBy"`   // 更新者
	UpdateTime time.Time `json:"updateTime"` // 更新时间
}

type DeptTreeByRoleId struct {
	CheckedKeys []int64           `json:"checkedKeys"`
	DeptList    []*system.SysDept `json:"deptList"`
}

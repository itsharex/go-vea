// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package system

const TableNameSysUserRole = "sys_user_role"

// SysUserRole mapped from table <sys_user_role>
type SysUserRole struct {
	UserID int64 `gorm:"column:user_id;type:bigint;primaryKey" json:"userId"` // 用户ID
	RoleID int64 `gorm:"column:role_id;type:bigint;primaryKey" json:"roleId"` // 角色ID
}

// TableName SysUserRole's table name
func (*SysUserRole) TableName() string {
	return TableNameSysUserRole
}
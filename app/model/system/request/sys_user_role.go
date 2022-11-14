package request

type SysUserRole struct {
	UserID  int64   `json:"userId"`  // 用户ID
	RoleID  int64   `json:"roleId"`  // 角色ID
	UserIds []int64 `json:"userIds"` //userIds
	RoleIds []int64 `json:"roleIds"` // roleIds
}

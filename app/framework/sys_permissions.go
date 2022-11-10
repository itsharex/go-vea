package framework

import (
	"fmt"
	"go-vea/app/model/system/response"
)

type SysPermissionsService struct{}

var SysPermissionSrv = new(SysPermissionsService)

// GetRolePermission 获取角色数据权限
func (p *SysPermissionsService) GetRolePermission(sysUserResp *response.SysUserResp) ([]string, error) {
	var rolePerms []string
	sysUser := sysUserResp.SysUser
	if sysUser.IsAdmin(sysUser.UserID) {
		rolePerms = append(rolePerms, "*:*:*")
	} else {
		rolePerms = append(rolePerms, "addAll")
	}
	return rolePerms, nil
}

// GetMenuPermission 获取菜单数据权限
func (p *SysPermissionsService) GetMenuPermission(sysUserResp *response.SysUserResp) ([]string, error) {
	var menuPerms []string
	sysUser := sysUserResp.SysUser
	if sysUser.IsAdmin(sysUser.UserID) {
		menuPerms = append(menuPerms, "*:*:*")
	} else {
		roleIds := sysUserResp.RoleIds
		if roleIds != nil && len(roleIds) > 0 {
			// 多角色设置permissions属性，以便数据权限匹配权限
			for i, roleId := range roleIds {
				fmt.Println(i, roleId)
				//Set<String> rolePerms = menuService.selectMenuPermsByRoleId(role.getRoleId());
				//role.setPermissions(rolePerms);
				//perms.addAll(rolePerms);
			}
		} else {
			menuPerms = append(menuPerms, "menuService.selectMenuPermsByUserId(user.getUserId())")
		}
	}
	return menuPerms, nil
}

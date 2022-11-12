package framework

import (
	"context"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
)

type SysPermissionsService struct{}

var SysPermissionSrv = new(SysPermissionsService)

// GetRolePermission 获取角色数据权限
func (*SysPermissionsService) GetRolePermission(ctx context.Context, sysUserResp *response.SysUserResp) ([]string, error) {
	var rolePerms []string
	sysUser := sysUserResp.SysUser
	if sysUser.IsAdmin(sysUser.UserID) {
		rolePerms = append(rolePerms, "*:*:*")
	} else {
		perms, err := syssrv.SysRoleSrv.SelectRolePermissionByUserId(ctx, sysUser)
		if err != nil {
			return nil, err
		}
		rolePerms = append(rolePerms, perms...)
	}
	return rolePerms, nil
}

// GetMenuPermission 获取菜单数据权限
func (*SysPermissionsService) GetMenuPermission(ctx context.Context, sysUserResp *response.SysUserResp) ([]string, error) {
	sysMenuDao := sysdao.NewSysMenuDao(ctx)
	var menuPerms []string
	sysUser := sysUserResp.SysUser
	if sysUser.IsAdmin(sysUser.UserID) {
		menuPerms = append(menuPerms, "*:*:*")
	} else {
		roles := sysUserResp.Roles
		if roles != nil {
			// 多角色设置permissions属性，以便数据权限匹配权限
			for _, role := range roles {
				rolePerms, err := sysMenuDao.SelectMenuPermsByRoleId(role.RoleID)
				if err != nil {
					return nil, err
				}
				role.Permissions = rolePerms
				menuPerms = append(menuPerms, rolePerms...)
			}
		} else {
			rolePerms, err := sysMenuDao.SelectMenuPermsByUserId(sysUserResp.SysUser.UserID)
			if err != nil {
				return nil, err
			}
			menuPerms = append(menuPerms, rolePerms...)
		}
	}
	return menuPerms, nil
}

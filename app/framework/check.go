package framework

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/service/syssrv"
	"go-vea/global"
)

type CheckService struct{}

var CheckSrv = new(CheckService)

// CheckRoleAllowed 校验角色是否允许操作
func (*CheckService) CheckRoleAllowed(ctx *gin.Context, sysRole *system.SysRole) bool {
	if sysRole != nil && sysRole.IsAdmin(sysRole.RoleID) {
		return false
	}
	return true
}

// CheckRoleDataScope 校验角色是否有数据权限
func (*CheckService) CheckRoleDataScope(ctx *gin.Context, roleId int64) bool {
	loginUser, _ := TokenSrv.GetLoginUser(ctx)
	if !loginUser.SysUserResp.SysUser.IsAdmin(loginUser.UserID) {
		params := &request.SysRole{
			RoleID: roleId,
		}
		data, _ := syssrv.SysRoleSrv.GetSysRoleList(ctx, params)
		if data.Rows == nil {
			global.Logger.Error("没有权限访问用户数据")
			return false
		}
	}
	return true
}

// CheckUserAllowed 校验用户是否允许操作
func (*CheckService) CheckUserAllowed(ctx *gin.Context, sysUser *system.SysUser) bool {
	loginUser, _ := TokenSrv.GetLoginUser(ctx)
	if sysUser != nil && loginUser.SysUserResp.SysUser.IsAdmin(sysUser.UserID) {
		global.Logger.Error("不允许操作超级管理员用户")
		return false
	}
	return true
}

// CheckUserDataScope 校验用户是否有数据权限
func (*CheckService) CheckUserDataScope(ctx *gin.Context, userId int64) bool {
	loginUser, _ := TokenSrv.GetLoginUser(ctx)
	if !loginUser.SysUserResp.SysUser.IsAdmin(userId) {
		params := &request.SysUser{
			UserID: userId,
		}
		data, _ := syssrv.SysUserSrv.GetSysUserList(ctx, params)
		if data.Rows == nil {
			global.Logger.Error("没有权限访问用户数据")
			return false
		}
	}
	return true
}

// CheckDeptDataScope 校验部门是否有数据权限
func (*CheckService) CheckDeptDataScope(ctx *gin.Context, deptId int64) bool {
	loginUser, _ := TokenSrv.GetLoginUser(ctx)
	if !loginUser.SysUserResp.SysUser.IsAdmin(loginUser.UserID) {
		params := &request.SysDept{
			DeptID: deptId,
		}
		data, _ := syssrv.SysDeptSrv.GetSysDeptList(ctx, params)
		if data.Rows == nil {
			global.Logger.Error("没有权限访问部门数据")
			return false
		}
	}
	return true
}

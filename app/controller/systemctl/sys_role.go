package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/framework"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysRoleApi struct{}

func (*SysRoleApi) GetSysRoleList(ctx *gin.Context) {
	var params request.SysRole
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysRoleSrv.GetSysRoleList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysRoleApi) GetSysRole(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))
	data, err := syssrv.SysRoleSrv.GetSysRoleById(ctx, int64(roleId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysRoleApi) AddSysRole(ctx *gin.Context) {
	var params system.SysRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysRoleSrv.AddSysRole(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) UpdateSysRole(ctx *gin.Context) {
	var params system.SysRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysRoleSrv.UpdateRoleById(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		// 更新缓存用户权限
		loginUser, _ := framework.TokenSrv.GetLoginUser(ctx)
		sysUser := loginUser.SysUserResp.SysUser
		if sysUser != nil && !sysUser.IsAdmin(sysUser.UserID) {
			perms, _ := framework.SysPermissionSrv.GetMenuPermission(ctx, loginUser.SysUserResp)
			userInfo, _ := syssrv.SysUserSrv.SelectUserByUsername(ctx, sysUser.Username)
			loginUser.Permissions = perms
			loginUser.SysUserResp.SysUser = userInfo
			framework.TokenSrv.SetLoginUser(loginUser)
		}
		result.Ok(ctx)
	}
}

func (*SysRoleApi) DeleteSysRole(ctx *gin.Context) {
	var params request.SysRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysRoleSrv.DeleteSysRoleByIds(ctx, params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) DataScope(ctx *gin.Context) {
	var params system.SysRole
	_ = ctx.ShouldBindJSON(&params)
	if !framework.CheckSrv.CheckRoleAllowed(ctx, &params) {
		result.FailWithMessage("不允许操作超级管理员用户", ctx)
		return
	}
	if !framework.CheckSrv.CheckRoleDataScope(ctx, params.RoleID) {
		result.FailWithMessage("没有权限访问角色数据", ctx)
		return
	}
	err := syssrv.SysRoleSrv.AuthDataScope(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) ChangeStatus(ctx *gin.Context) {
	var params system.SysRole
	_ = ctx.ShouldBindJSON(&params)
	if !framework.CheckSrv.CheckRoleAllowed(ctx, &params) {
		result.FailWithMessage("不允许操作超级管理员用户", ctx)
		return
	}
	if !framework.CheckSrv.CheckRoleDataScope(ctx, params.RoleID) {
		result.FailWithMessage("没有权限访问角色数据", ctx)
		return
	}
	err := syssrv.SysRoleSrv.UpdateRoleStatus(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) OptionSelect(ctx *gin.Context) {
	data, err := syssrv.SysRoleSrv.SelectRoleAll(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysRoleApi) AllocatedList(ctx *gin.Context) {
	var params request.SysUser
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysUserSrv.SelectAllocatedList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysRoleApi) UnallocatedList(ctx *gin.Context) {
	var params request.SysUser
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysUserSrv.SelectUnallocatedList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysRoleApi) CancelAuthUser(ctx *gin.Context) {
	var params system.SysUserRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysUserSrv.DeleteAuthUser(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) CancelAuthUserAll(ctx *gin.Context) {
	var params request.SysUserRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysUserSrv.BatchDeleteAuthUser(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) SelectAuthUserAll(ctx *gin.Context) {
	var params request.SysUserRole
	_ = ctx.ShouldBindJSON(&params)
	if !framework.CheckSrv.CheckRoleDataScope(ctx, params.RoleID) {
		result.FailWithMessage("没有权限访问角色数据", ctx)
		return
	}
	err := syssrv.SysUserSrv.BatchAddAuthUser(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) DeptTree(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))
	checkedKeys, err := syssrv.SysDeptSrv.SelectDeptListByRoleId(ctx, int64(roleId))
	deptList, err := syssrv.SysDeptSrv.GetDeptTreeList(ctx, &request.SysDept{})
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		res := &response.DeptTreeByRoleId{
			CheckedKeys: checkedKeys,
			DeptList:    deptList,
		}
		result.OkWithData(res, ctx)
	}
}

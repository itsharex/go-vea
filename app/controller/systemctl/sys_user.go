package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/framework"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
	"go-vea/util"
	"strconv"
)

type SysUserApi struct{}

func (*SysUserApi) GetSysUserList(ctx *gin.Context) {
	var params request.SysUser
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysUserSrv.GetSysUserList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysUserApi) GetSysUser(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))
	userInfo := &response.UserInfoById{}
	if !framework.CheckSrv.CheckUserDataScope(ctx, int64(userId)) {
		result.FailWithMessage("没有权限访问用户数据", ctx)
		return
	}
	var roles []*system.SysRole
	tmp, err := syssrv.SysRoleSrv.SelectRoleAll(ctx)
	for _, role := range tmp {
		if !role.IsAdmin(role.RoleID) {
			roles = append(roles, role)
		}
	}
	posts, err := syssrv.SysPostSrv.SelectPostAll(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
		return
	}
	// 角色选项
	userInfo.Roles = roles
	// 部门选项
	userInfo.Posts = posts
	if userId != -1 {
		data, err := syssrv.SysUserSrv.GetSysUserById(ctx, int64(userId))
		if err != nil {
			result.FailWithMessage(err.Error(), ctx)
			return
		} else {
			postIds, _ := syssrv.SysPostSrv.SelectPostListByUserId(ctx, int64(userId))
			// 用户信息
			userInfo.User = data
			// 已选择部门
			userInfo.PostIds = postIds
			// 已选择角色
			userInfo.RoleIds = data.RoleIds
		}
	}
	result.OkWithData(userInfo, ctx)
}

func (*SysUserApi) AddSysUser(ctx *gin.Context) {
	var params request.AddSysUser
	_ = ctx.ShouldBindJSON(&params)
	sysUser := params.SysUser
	if !syssrv.SysUserSrv.CheckUserNameUnique(ctx, sysUser) {
		result.FailWithMessage("新增用户'"+sysUser.Username+"'失败，登录账号已存在", ctx)
		return
	} else if sysUser.PhoneNumber != "" && !syssrv.SysUserSrv.CheckPhoneUnique(ctx, sysUser) {
		result.FailWithMessage("新增用户'"+sysUser.Username+"'失败，手机号码已存在", ctx)
		return
	} else if sysUser.Email != "" && !syssrv.SysUserSrv.CheckEmailUnique(ctx, sysUser) {
		result.FailWithMessage("新增用户'"+sysUser.Username+"'失败，邮箱账号已存在", ctx)
		return
	}
	err := syssrv.SysUserSrv.AddSysUser(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) UpdateSysUser(ctx *gin.Context) {
	var params request.AddSysUser
	_ = ctx.ShouldBindJSON(&params)
	sysUser := params.SysUser
	if !framework.CheckSrv.CheckUserAllowed(ctx, sysUser) {
		result.FailWithMessage("不允许操作超级管理员用户", ctx)
		return
	}
	if !framework.CheckSrv.CheckUserDataScope(ctx, sysUser.UserID) {
		result.FailWithMessage("没有权限访问用户数据", ctx)
		return
	}
	if !syssrv.SysUserSrv.CheckUserNameUnique(ctx, sysUser) {
		result.FailWithMessage("修改用户'"+sysUser.Username+"'失败，邮箱账号已存在", ctx)
		return
	} else if sysUser.PhoneNumber != "" && !syssrv.SysUserSrv.CheckPhoneUnique(ctx, sysUser) {
		result.FailWithMessage("修改用户'"+sysUser.Username+"'失败，邮箱账号已存在", ctx)
		return
	} else if sysUser.Email != "" && !syssrv.SysUserSrv.CheckEmailUnique(ctx, sysUser) {
		result.FailWithMessage("修改用户'"+sysUser.Username+"'失败，邮箱账号已存在", ctx)
		return
	}
	err := syssrv.SysUserSrv.UpdateUserById(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) DeleteSysUser(ctx *gin.Context) {
	var params request.SysUser
	_ = ctx.ShouldBindJSON(&params)
	loginUser, _ := framework.TokenSrv.GetLoginUser(ctx)
	var hasCurrentUserId bool
	for _, id := range params.Ids {
		if id == loginUser.UserID {
			hasCurrentUserId = true
			break
		}
	}
	if hasCurrentUserId {
		result.FailWithMessage("当前用户不能删除", ctx)
		return
	}
	err := syssrv.SysUserSrv.DeleteSysUserByIds(ctx, params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) ResetPwd(ctx *gin.Context) {
	var params system.SysUser
	_ = ctx.ShouldBindJSON(&params)
	if !framework.CheckSrv.CheckUserAllowed(ctx, &params) {
		result.FailWithMessage("不允许操作超级管理员用户", ctx)
		return
	}
	if !framework.CheckSrv.CheckUserDataScope(ctx, params.UserID) {
		result.FailWithMessage("没有权限访问用户数据", ctx)
		return
	}
	loginUser, _ := framework.TokenSrv.GetLoginUser(ctx)
	params.UpdateBy = loginUser.SysUserResp.SysUser.Username
	params.Password, _ = util.PasswordHash(params.Password)
	err := syssrv.SysUserSrv.ResetPwd(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) ChangeStatus(ctx *gin.Context) {
	var params system.SysUser
	_ = ctx.ShouldBindJSON(&params)
	if !framework.CheckSrv.CheckUserAllowed(ctx, &params) {
		result.FailWithMessage("不允许操作超级管理员用户", ctx)
		return
	}
	if !framework.CheckSrv.CheckUserDataScope(ctx, params.UserID) {
		result.FailWithMessage("没有权限访问用户数据", ctx)
		return
	}
	loginUser, _ := framework.TokenSrv.GetLoginUser(ctx)
	params.UpdateBy = loginUser.SysUserResp.SysUser.Username
	err := syssrv.SysUserSrv.UpdateUserStatus(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) AuthRole(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))
	userInfo := response.UserInfoById{}

	data, err := syssrv.SysUserSrv.GetSysUserById(ctx, int64(userId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		var roles []*system.SysRole
		allRoles, _ := syssrv.SysRoleSrv.SelectRolesByUserId(ctx, int64(userId))
		var sysUser *system.SysUser
		if sysUser.IsAdmin(int64(userId)) {
			userInfo.Roles = allRoles
		} else {
			for _, role := range allRoles {
				if !role.IsAdmin(int64(userId)) {
					roles = append(roles, role)
				}
			}
			userInfo.Roles = roles
		}
		userInfo.User = data
		result.OkWithData(userInfo, ctx)
	}
}

func (*SysUserApi) InsertAuthRole(ctx *gin.Context) {
	var params request.AddUserRole
	_ = ctx.ShouldBindJSON(&params)
	if !framework.CheckSrv.CheckUserDataScope(ctx, params.UserId) {
		result.FailWithMessage("没有权限访问用户数据", ctx)
		return
	}
	err := syssrv.SysUserSrv.InsertUserAuth(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) GetDeptTree(ctx *gin.Context) {
	params := request.SysDept{}
	data, err := syssrv.SysDeptSrv.GetDeptTreeList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

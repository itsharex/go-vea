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
)

type SysProfileApi struct{}

func (*SysProfileApi) GetProfile(ctx *gin.Context) {
	loginUser, err := framework.TokenSrv.GetLoginUser(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		roleGroup, _ := syssrv.SysUserSrv.SelectUserRoleGroup(ctx, loginUser.SysUserResp.SysUser.Username)
		postGroup, _ := syssrv.SysUserSrv.SelectUserPostGroup(ctx, loginUser.SysUserResp.SysUser.Username)
		userProfile := response.UserProfile{
			UserInfo:  loginUser.SysUserResp.SysUser,
			SysDept:   loginUser.SysUserResp.SysDept,
			RoleGroup: roleGroup,
			PostGroup: postGroup,
		}
		result.OkWithData(userProfile, ctx)
	}
}

func (*SysProfileApi) UpdateProfile(ctx *gin.Context) {
	var params system.SysUser
	_ = ctx.ShouldBindJSON(&params)
	loginUser, err := framework.TokenSrv.GetLoginUser(ctx)
	phoneNumberUnique := syssrv.SysUserSrv.CheckPhoneUnique(ctx, &params)
	if params.PhoneNumber != "" && !phoneNumberUnique {
		result.FailWithMessage("修改用户'"+params.Username+"'失败，手机号码已存在", ctx)
		return
	}
	emailUnique := syssrv.SysUserSrv.CheckEmailUnique(ctx, &params)
	if params.Email != "" && !emailUnique {
		result.FailWithMessage("修改用户'"+params.Username+"'失败，邮箱账号已存在", ctx)
		return
	}
	params.UserID = loginUser.UserID
	err = syssrv.SysUserSrv.UpdateUserProfile(ctx, &params)
	if err != nil {
		result.FailWithMessage("修改个人信息异常，请联系管理员", ctx)
	} else {
		// 更新缓存
		loginUser.SysUserResp.SysUser.Nickname = params.Nickname
		loginUser.SysUserResp.SysUser.PhoneNumber = params.PhoneNumber
		loginUser.SysUserResp.SysUser.Email = params.Email
		loginUser.SysUserResp.SysUser.Gender = params.Gender
		framework.TokenSrv.SetLoginUser(loginUser)
		result.Ok(ctx)
	}
}

func (*SysProfileApi) UpdatePassword(ctx *gin.Context) {
	var params request.ResetPwd
	_ = ctx.ShouldBindJSON(&params)
	loginUser, err := framework.TokenSrv.GetLoginUser(ctx)
	username := loginUser.SysUserResp.SysUser.Username
	hashPassword := loginUser.SysUserResp.SysUser.Password
	params.Username = username
	if !util.PasswordVerify(params.OldPassword, hashPassword) {
		result.FailWithMessage("修改密码失败，旧密码错误", ctx)
		return
	}
	if util.PasswordVerify(params.NewPassword, hashPassword) {
		result.FailWithMessage("新密码不能与旧密码相同", ctx)
		return
	}
	params.NewPassword, _ = util.PasswordHash(params.NewPassword)
	err = syssrv.SysUserSrv.ResetUserPwd(ctx, &params)
	if err != nil {
		result.FailWithMessage("修改密码异常，请联系管理员", ctx)
	} else {
		// 更新缓存用户密码
		loginUser.SysUserResp.SysUser.Password = params.NewPassword
		framework.TokenSrv.SetLoginUser(loginUser)
		result.Ok(ctx)
	}
}

func (*SysProfileApi) UploadAvatar(ctx *gin.Context) {
	result.Ok(ctx)
}

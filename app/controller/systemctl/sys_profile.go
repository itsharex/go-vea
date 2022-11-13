package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/framework"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
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
	result.Ok(ctx)
}

func (*SysProfileApi) UpdatePassword(ctx *gin.Context) {
	result.Ok(ctx)
}

func (*SysProfileApi) UploadAvatar(ctx *gin.Context) {
	result.Ok(ctx)
}

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
		roleGroup := syssrv.SysUserSrv.SelectUserRoleGroup(ctx, loginUser.SysUserResp.SysUser.UserName)
		postGroup := syssrv.SysUserSrv.SelectUserPostGroup(ctx, loginUser.SysUserResp.SysUser.UserName)
		userProfile := response.UserProfile{
			UserInfo:  loginUser.SysUserResp.SysUser,
			RoleGroup: roleGroup,
			PostGroup: postGroup,
		}
		result.OkWithData(userProfile, ctx)
	}
}

func (*SysProfileApi) UpdateProfile(ctx *gin.Context) {

}

func (*SysProfileApi) UpdatePassword(ctx *gin.Context) {

}

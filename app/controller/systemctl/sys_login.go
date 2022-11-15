package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/framework"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
	"go-vea/global"
)

type SysLoginApi struct{}

// Login 用户登录
func (s *SysLoginApi) Login(ctx *gin.Context) {
	var loginBody request.LoginBody
	_ = ctx.ShouldBindJSON(&loginBody)
	token, err := framework.SysLoginSrv.Login(ctx, &loginBody)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithMessage(token, ctx)
	}
}

// Logout 退出登录
func (s *SysLoginApi) Logout(ctx *gin.Context) {
	err := framework.TokenSrv.Logout(ctx)
	if err != nil {
		result.Fail(ctx)
	} else {
		result.Ok(ctx)
	}
}

// GetUserInfo 获取当前登录用户信息
func (s *SysLoginApi) GetUserInfo(ctx *gin.Context) {
	loginUser, err := framework.TokenSrv.GetLoginUser(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		roles, _ := syssrv.SysRoleSrv.GetRolePermission(ctx, loginUser.SysUserResp.SysUser)
		userInfo := &response.UserInfo{
			User:        loginUser.SysUserResp.SysUser,
			Roles:       roles,
			Permissions: loginUser.Permissions,
		}
		result.OkWithData(userInfo, ctx)
	}
}

// GetRouters 获取前端路由信息 菜单
func (s *SysLoginApi) GetRouters(ctx *gin.Context) {
	// 菜单树
	loginUser, err := framework.TokenSrv.GetLoginUser(ctx)
	if err != nil {
		global.Logger.Error(err)
	}
	menus, err := syssrv.SysMenuSrv.SelectMenuTreeByUserId(ctx, loginUser.SysUserResp.SysUser)
	if err != nil {
		result.Fail(ctx)
	} else {
		result.OkWithData(syssrv.SysMenuSrv.GetBuildMenus(menus), ctx)
	}
}

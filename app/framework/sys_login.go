package framework

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-web-template/app/dao/sysdao"
	"go-web-template/app/model/system"
	"go-web-template/app/model/system/request"
	"go-web-template/app/model/system/response"
	"go-web-template/app/service/syssrv"
	"go-web-template/global"
	"go-web-template/util/captcha"
)

type SysLoginService struct{}

var SysLoginSrv = new(SysLoginService)

func (s *SysLoginService) Login(ctx *gin.Context, loginBody *request.LoginBody) (string, error) {
	captchaEnabled, _ := syssrv.SysConfigSrv.SelectCaptchaEnabled(ctx)
	// 验证码开关
	if captchaEnabled {
		// 验证码验证
		r := captcha.CaptVerify(loginBody.UUID, loginBody.Code)
		if !r {
			err := errors.New("验证码错误")
			global.Logger.Error("验证码错误", err)
			return "", err
		}
	}
	loginUser, err := loadUserByUsername(ctx, loginBody)
	if err != nil {
		global.Logger.Error(err)
		return "", err
	}
	token, err := TokenSrv.CreateToken(loginUser)
	if err != nil {
		global.Logger.Error(err)
		return "", err
	}
	recordLoginInfo(loginUser.UserID)
	return token, err
}

func recordLoginInfo(userId int64) {
	// todo 记录登录信息
}

func loadUserByUsername(ctx *gin.Context, loginBody *request.LoginBody) (*response.LoginUser, error) {
	userDao := sysdao.NewSysUserDao(ctx)
	sysUser, err := userDao.SelectUserByUserName(loginBody.UserName)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	// 账号验证
	if sysUser == nil {
		global.Logger.Error("用户不存在")
		return nil, errors.New("用户不存在")
	} else if sysUser.DelFlag == "2" {
		global.Logger.Error("账号已删除")
		return nil, errors.New("账号已删除")
	} else if sysUser.Status == "1" {
		global.Logger.Error("账号已停用")
		return nil, errors.New("账号已停用")
	}
	// 密码验证
	err = SysPasswordSrv.Validate(sysUser, loginBody)
	if err != nil {
		global.Logger.Error("密码错误", err)
		return nil, err
	}
	sysUserResp := buildSysUserResp(sysUser)
	return createLoginUser(sysUserResp), err
}

func buildSysUserResp(sysUser *system.SysUser) *response.SysUserResp {
	// todo
	roles := []string{"admin"}
	roleIds := []int64{1}
	postIds := []int64{1}
	sysUserResp := &response.SysUserResp{
		SysUser: sysUser,
		SysDept: nil,
		Roles:   roles,
		RoleIds: roleIds,
		RoleId:  1,
		PostIds: postIds,
	}
	return sysUserResp
}

func createLoginUser(sysUserResp *response.SysUserResp) *response.LoginUser {
	permissions, err := SysPermissionSrv.GetMenuPermission(sysUserResp)
	if err != nil {
		global.Logger.Error(err)
	}
	loginUser := response.LoginUser{
		UserID:      sysUserResp.SysUser.UserID,
		DeptID:      sysUserResp.SysUser.DeptID,
		Permissions: permissions,
		SysUserResp: sysUserResp,
	}
	return &loginUser
}

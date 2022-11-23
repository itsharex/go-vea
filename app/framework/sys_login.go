package framework

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
	"go-vea/global"
	"go-vea/util/captcha"
	"time"
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
	token, err := TokenSrv.CreateToken(ctx, loginUser)
	if err != nil {
		global.Logger.Error(err)
		return "", err
	}
	err = recordLoginInfo(ctx, loginUser.UserID)
	return token, err
}

func recordLoginInfo(ctx *gin.Context, userId int64) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	user := &system.SysUser{}
	user.UserID = userId
	// todo
	user.LoginIP = "127.0.0.1"
	now := time.Now()
	user.LoginDate = &now

	err := sysUserDao.UpdateUserProfile(user)
	if err != nil {
		return err
	}
	return nil
}

func loadUserByUsername(ctx *gin.Context, loginBody *request.LoginBody) (*response.LoginUser, error) {
	userDao := sysdao.NewSysUserDao(ctx)
	sysUser, err := userDao.SelectUserByUsername(loginBody.Username)
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
	sysUserResp, _ := buildSysUserResp(ctx, sysUser)
	return createLoginUser(ctx, sysUserResp), err
}

func buildSysUserResp(ctx context.Context, sysUser *system.SysUser) (*response.SysUserResp, error) {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	sysUserRoleDao := sysdao.NewSysUserRoleDao(ctx)
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	var dept *system.SysDept
	var roles []*system.SysRole
	var roleIds []int64
	sysUserResp := &response.SysUserResp{}
	var err error
	dept, err = sysDeptDao.SelectById(sysUser.DeptID)
	roleIds, err = sysUserRoleDao.SelectByUserId(sysUser.UserID)
	if roleIds != nil {
		roles, err = sysRoleDao.SelectByRoleIds(roleIds)
		sysUserResp.Roles = roles
		sysUserResp.RoleIds = roleIds
	}
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	sysUserResp.SysUser = sysUser
	sysUserResp.SysDept = dept
	return sysUserResp, err
}

func createLoginUser(ctx context.Context, sysUserResp *response.SysUserResp) *response.LoginUser {
	permissions, err := SysPermissionSrv.GetMenuPermission(ctx, sysUserResp)
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

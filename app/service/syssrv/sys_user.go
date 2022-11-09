package syssrv

import (
	"context"
	"go-web-template/app/common/page"
	"go-web-template/app/dao/sysdao"
	"go-web-template/app/model/system"
	"go-web-template/app/model/system/request"
	"go-web-template/global"
)

type SysUserService struct{}

var SysUserSrv = new(SysUserService)

func (*SysUserService) GetSysUserList(ctx context.Context, sysUser *request.SysUser) (*page.Pagination, error) {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	data, err := sysUserDao.SelectList(sysUser)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysUserService) GetSysUserById(ctx context.Context, roleId int64) (*system.SysUser, error) {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	data, err := sysUserDao.SelectById(roleId)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysUserService) AddSysUser(ctx context.Context, sysUser *system.SysUser) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	err := sysUserDao.Insert(sysUser)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysUserService) UpdateUserById(ctx context.Context, sysUser *system.SysUser) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	err := sysUserDao.UpdateById(sysUser)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysUserService) DeleteSysUserById(ctx context.Context, userId int64) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	err := sysUserDao.DeleteById(userId)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	// 删除用户与角色关联 todo
	// 删除用户与岗位关联 todo
	return nil
}

func (*SysUserService) DeleteSysUserByIds(ctx context.Context, ids []int64) error {
	sysUserDao := sysdao.NewSysUserDao(ctx)
	err := sysUserDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	// 删除用户与角色关联 todo
	// 删除用户与岗位关联 todo
	return nil
}

func (*SysUserService) SelectUserRoleGroup(ctx context.Context, name string) string {
	return ""
}

func (*SysUserService) SelectUserPostGroup(ctx context.Context, name string) string {
	return ""
}

func checkUserNameUnique(ctx context.Context, sysUser *system.SysUser) bool {
	return false
}

func checkPhoneUnique(ctx context.Context, sysUser *system.SysUser) bool {
	return false
}

func checkEmailUnique(ctx context.Context, sysUser *system.SysUser) bool {
	return false
}

func checkUserAllowed(ctx context.Context, sysUser *system.SysUser) bool {
	return false
}

func checkUserDataScope(ctx context.Context, sysUser *system.SysUser) bool {
	return false
}

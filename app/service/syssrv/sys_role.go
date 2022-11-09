package syssrv

import (
	"context"
	"errors"
	"go-web-template/app/common/page"
	"go-web-template/app/dao/sysdao"
	"go-web-template/app/model/system"
	"go-web-template/app/model/system/request"
	"go-web-template/global"
)

type SysRoleService struct{}

var SysRoleSrv = new(SysRoleService)

func (*SysRoleService) GetSysRoleList(ctx context.Context, sysRole *request.SysRole) (*page.Pagination, error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	data, err := sysRoleDao.SelectList(sysRole)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysRoleService) GetSysRoleById(ctx context.Context, roleId int64) (*system.SysRole, error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	data, err := sysRoleDao.SelectById(roleId)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysRoleService) AddSysRole(ctx context.Context, sysRole *system.SysRole) error {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	hasRoleName := checkRoleNameUnique(ctx, sysRole)
	hasRoleKey := checkRoleKeyUnique(ctx, sysRole)
	if hasRoleName {
		global.Logger.Error("新增失败！角色名称已存在")
		return errors.New("新增失败！角色名称已存在")
	} else if hasRoleKey {
		global.Logger.Error("新增失败！角色key已存在")
		return errors.New("新增失败！角色key已存在")
	}
	err := sysRoleDao.Insert(sysRole)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysRoleService) UpdateRoleById(ctx context.Context, sysRole *system.SysRole) error {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	hasRoleName := checkRoleNameUnique(ctx, sysRole)
	hasRoleKey := checkRoleKeyUnique(ctx, sysRole)
	if hasRoleName {
		global.Logger.Error("修改失败！角色名称已存在")
		return errors.New("修改失败！角色名称已存在")
	} else if hasRoleKey {
		global.Logger.Error("修改失败！角色key已存在")
		return errors.New("修改失败！角色key已存在")
	}
	err := sysRoleDao.UpdateById(sysRole)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysRoleService) DeleteSysRoleByIds(ctx context.Context, ids []int64) error {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	err := sysRoleDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func checkRoleNameUnique(ctx context.Context, sysRole *system.SysRole) bool {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	r, e := sysRoleDao.CheckRoleNameUnique(sysRole.RoleName)
	if e != nil {
		global.Logger.Error(e)
		return false
	}
	return r > 0
}

func checkRoleKeyUnique(ctx context.Context, sysRole *system.SysRole) bool {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	r, e := sysRoleDao.CheckRoleKeyUnique(sysRole.RoleName)
	if e != nil {
		global.Logger.Error(e)
		return false
	}
	return r > 0
}

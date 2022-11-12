package syssrv

import (
	"context"
	"errors"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
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

func (*SysRoleService) SelectRolesByUserId(ctx context.Context, userId int64) ([]*system.SysRole, error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	userRoles, err := sysRoleDao.SelectRolePermissionByUserId(userId)
	roles, err := sysRoleDao.SelectAll(&request.SysRole{})
	for i, role := range roles {
		for _, userRole := range userRoles {
			if role.RoleID == userRole.RoleID {
				roles[i].Flag = true
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *SysRoleService) SelectRollAll(ctx context.Context) ([]*system.SysRole, error) {
	data, err := s.GetSysRoleList(ctx, &request.SysRole{})
	if err != nil {
		return nil, err
	}
	roles, ok := data.Rows.([]*system.SysRole)
	if !ok {
		global.Logger.Error("类型转换错误")
		return nil, errors.New("类型转换错误")
	}
	return roles, err
}

func (*SysRoleService) GetRolePermission(ctx context.Context, sysUser *system.SysUser) ([]string, error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	var perms []string
	if sysUser.IsAdmin(sysUser.UserID) {
		perms = append(perms, "admin")
	} else {
		roles, err := sysRoleDao.SelectRolePermissionByUserId(sysUser.UserID)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			perms = append(perms, role.RoleKey)
		}
	}
	return perms, nil
}

func (*SysRoleService) SelectRolePermissionByUserId(ctx context.Context, user *system.SysUser) ([]string, error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	var perms []string
	roles, err := sysRoleDao.SelectRolePermissionByUserId(user.UserID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		perms = append(perms, role.RoleKey)
	}
	return perms, nil
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

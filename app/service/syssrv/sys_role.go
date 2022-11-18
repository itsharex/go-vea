package syssrv

import (
	"context"
	"errors"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
	"gorm.io/gorm"
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

func (*SysRoleService) AddSysRole(ctx context.Context, sysRole *system.SysRole) (err error) {
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
	// todo Transactional
	// 新增角色信息
	err = sysRoleDao.Insert(sysRole)
	// 新增角色菜单关联
	err = insertRoleMenu(ctx, sysRole)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysRoleService) UpdateRoleById(ctx context.Context, sysRole *system.SysRole) (err error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	sysRoleMenuDao := sysdao.NewSysRoleMenuDao(ctx)
	// todo Transactional
	hasRoleName := checkRoleNameUnique(ctx, sysRole)
	hasRoleKey := checkRoleKeyUnique(ctx, sysRole)
	if hasRoleName {
		global.Logger.Error("修改失败！角色名称已存在")
		return errors.New("修改失败！角色名称已存在")
	} else if hasRoleKey {
		global.Logger.Error("修改失败！角色key已存在")
		return errors.New("修改失败！角色key已存在")
	}
	// 修改角色信息
	err = sysRoleDao.UpdateById(sysRole)
	// 删除角色与菜单关联
	err = sysRoleMenuDao.DeleteRoleMenuByRoleId(sysRole.RoleID)
	// 新增角色菜单关联
	err = insertRoleMenu(ctx, sysRole)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysRoleService) DeleteSysRoleByIds(ctx context.Context, ids []int64) (err error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	for _, roleId := range ids {
		used := countUserRoleByRoleId(ctx, roleId)
		if used {
			return errors.New("该角色已分配，不能删除")
		}
	}
	// todo Transactional
	// 删除角色与菜单关联
	roleMenuDao := sysdao.NewSysRoleMenuDao(ctx)
	err = roleMenuDao.DeleteRoleMenu(ids)
	// 删除角色与部门关联
	roleDeptDao := sysdao.NewSysRoleDeptDao(ctx)
	err = roleDeptDao.DeleteRoleDept(ids)
	// 删除角色信息
	err = sysRoleDao.DeleteByIds(ids)
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

func (s *SysRoleService) SelectRoleAll(ctx context.Context) ([]*system.SysRole, error) {
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

func (*SysRoleService) AuthDataScope(ctx context.Context, sysRole *system.SysRole) (err error) {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	sysRoleDeptDao := sysdao.NewSysRoleDeptDao(ctx)
	// todo Transactional
	// 修改角色信息
	err = sysRoleDao.UpdateById(sysRole)
	// 删除角色与部门关联
	err = sysRoleDeptDao.DeleteRoleDeptByRoleId(sysRole.RoleID)
	// 新增角色和部门信息（数据权限）
	err = insertRoleDept(ctx, sysRole)
	if err != nil {
		return err
	}
	return nil
}

func (*SysRoleService) UpdateRoleStatus(ctx context.Context, sysRole *system.SysRole) error {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	err := sysRoleDao.UpdateById(sysRole)
	if err != nil {
		return err
	}
	return nil
}

func checkRoleNameUnique(ctx context.Context, sysRole *system.SysRole) bool {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	data, err := sysRoleDao.CheckRoleNameUnique(sysRole.RoleName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return true
	}
	if data.RoleID != 0 && data.RoleID != sysRole.RoleID {
		return true
	}
	return false
}

func checkRoleKeyUnique(ctx context.Context, sysRole *system.SysRole) bool {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	data, err := sysRoleDao.CheckRoleKeyUnique(sysRole.RoleKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return true
	}
	if data.RoleID != 0 && data.RoleID != sysRole.RoleID {
		return true
	}
	return false
}

func countUserRoleByRoleId(ctx context.Context, roleId int64) bool {
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	count, err := sysRoleDao.CountUserRoleByRoleId(roleId)
	if err != nil {
		return true
	}
	return count > 0
}

func insertRoleMenu(ctx context.Context, sysRole *system.SysRole) error {
	sysRoleMenuDao := sysdao.NewSysRoleMenuDao(ctx)
	var sysRoleMenuList []*system.SysRoleMenu
	for _, menuId := range sysRole.MenuIds {
		sysRoleMenu := &system.SysRoleMenu{
			RoleID: sysRole.RoleID,
			MenuID: menuId,
		}
		sysRoleMenuList = append(sysRoleMenuList, sysRoleMenu)
	}
	if len(sysRoleMenuList) > 0 {
		err := sysRoleMenuDao.BatchRoleMenu(sysRoleMenuList)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertRoleDept(ctx context.Context, sysRole *system.SysRole) error {
	sysRoleDeptDao := sysdao.NewSysRoleDeptDao(ctx)
	var sysRoleDeptList []*system.SysRoleDept
	for _, deptId := range sysRole.DeptIds {
		sysRoleDept := &system.SysRoleDept{
			RoleID: sysRole.RoleID,
			DeptID: deptId,
		}
		sysRoleDeptList = append(sysRoleDeptList, sysRoleDept)
	}
	if len(sysRoleDeptList) > 0 {
		err := sysRoleDeptDao.BatchRoleDept(sysRoleDeptList)
		if err != nil {
			return err
		}
	}
	return nil
}

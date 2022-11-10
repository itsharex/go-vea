package sysdao

import (
	"context"
	"go-vea/app/common/e"
	"go-vea/app/common/page"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/configs"
	"gorm.io/gorm"
)

type SysRoleDao struct {
	*gorm.DB
}

func NewSysRoleDao(ctx context.Context) *SysRoleDao {
	return &SysRoleDao{configs.GetDB(ctx)}
}

func NewSysRoleDaoByDB(db *gorm.DB) *SysRoleDao {
	return &SysRoleDao{db}
}

func (dao *SysRoleDao) SelectList(sysRole *request.SysRole) (p *page.Pagination, err error) {
	var roleList []*system.SysRole
	p = new(page.Pagination)

	if sysRole.RoleID != 0 {
		dao.DB = dao.DB.Where("role_id = ?", sysRole.RoleID)
	}
	if sysRole.RoleName != "" {
		dao.DB = dao.DB.Where("role_name = ?", sysRole.RoleName)
	}
	if sysRole.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysRole.Status)
	}
	if sysRole.RoleKey != "" {
		dao.DB = dao.DB.Where("role_key = ?", sysRole.RoleKey)
	}

	if sysRole.OpenPage {
		p.PageNum = sysRole.PageNum
		p.PageSize = sysRole.PageSize
		err = dao.DB.Scopes(page.SelectPage(roleList, p, dao.DB)).Find(&roleList).Error
	} else {
		err = dao.DB.Find(&roleList).Error
	}
	p.Rows = roleList
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}

func (dao *SysRoleDao) SelectAll(sysRole *request.SysRole) (list []system.SysRole, err error) {
	if sysRole.RoleID != 0 {
		dao.DB = dao.DB.Where("role_id = ?", sysRole.RoleID)
	}
	if sysRole.RoleName != "" {
		dao.DB = dao.DB.Where("role_name = ?", sysRole.RoleName)
	}
	if sysRole.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysRole.Status)
	}
	if sysRole.RoleKey != "" {
		dao.DB = dao.DB.Where("role_key = ?", sysRole.RoleKey)
	}
	if sysRole.DataScope != "" {
		// todo
	}

	err = dao.DB.Where("del_flag = '0'").Find(&list).Error
	return
}

func (dao *SysRoleDao) SelectById(id int64) (sysRole *system.SysRole, err error) {
	err = dao.DB.Model(&sysRole).Where("role_id=?", id).Find(&sysRole).Error
	return
}

func (dao *SysRoleDao) Insert(sysRole *system.SysRole) error {
	return dao.DB.Model(&system.SysRole{}).Create(sysRole).Error
}

func (dao *SysRoleDao) UpdateById(sysRole *system.SysRole) error {
	// Save 保存更新数据库中的值。如果值不包含一个匹配的主键，值将被插入。
	return dao.DB.Save(sysRole).Error
}

func (dao *SysRoleDao) DeleteById(id int64) error {
	return dao.DB.Where("role_id = ?", id).Delete(&system.SysRole{}).Error
}

func (dao *SysRoleDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("role_id in (?)", ids).Delete(&system.SysRole{}).Error
}

func (dao *SysRoleDao) CheckRoleNameUnique(roleName string) (count int64, err error) {
	err = dao.DB.Model(&system.SysRole{}).Where("role_name = ?", roleName).Where("del_flag = '0'").Count(&count).Error
	return
}

func (dao *SysRoleDao) CheckRoleKeyUnique(roleKey string) (count int64, err error) {
	err = dao.DB.Model(&system.SysRole{}).Where("role_key = ?", roleKey).Where("del_flag = '0'").Count(&count).Error
	return
}

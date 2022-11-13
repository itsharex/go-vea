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

func (dao *SysRoleDao) SelectAll(sysRole *request.SysRole) (list []*system.SysRole, err error) {
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

func (dao *SysRoleDao) SelectRolePermissionByUserId(userId int64) (sysRoles []system.SysRole, err error) {
	err = dao.DB.Table("sys_role r").
		Distinct("r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, r.dept_check_strictly,r.status, r.del_flag, r.create_time, r.remark").
		Select("r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, r.dept_check_strictly,r.status, r.del_flag, r.create_time, r.remark").
		Joins("left join sys_user_role ur on ur.role_id = r.role_id").
		Joins("left join sys_user u on u.user_id = ur.user_id").
		Joins("left join sys_dept d on u.dept_id = d.dept_id").
		Where("r.del_flag = '0' and ur.user_id = ?", userId).
		Find(&sysRoles).Error
	return
}

func (dao *SysRoleDao) SelectByRoleIds(ids []int64) (roles []*system.SysRole, err error) {
	err = dao.DB.Where("role_id in (?)", ids).Find(&roles).Error
	return
}

func (dao *SysRoleDao) SelectRolesByUserName(username string) (roles []*system.SysRole, err error) {
	err = dao.DB.Table("sys_role r").
		Distinct("r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, r.dept_check_strictly,r.status, r.del_flag, r.create_time, r.remark ").
		Joins("left join sys_user_role ur on ur.role_id = r.role_id").
		Joins("left join sys_user u on u.user_id = ur.user_id").
		Joins("left join sys_dept d on u.dept_id = d.dept_id").
		Where("r.del_flag ='0' and u.username = ?", username).
		Find(&roles).Error
	return
}

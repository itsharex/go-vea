package sysdao

import (
	"context"
	"go-vea/app/model/system"
	"go-vea/configs"
	"gorm.io/gorm"
)

type SysRoleMenuDao struct {
	*gorm.DB
}

func NewSysRoleMenuDao(ctx context.Context) *SysRoleMenuDao {
	return &SysRoleMenuDao{configs.GetDB(ctx)}
}

func NewSysRoleMenuDaoByDB(db *gorm.DB) *SysRoleMenuDao {
	return &SysRoleMenuDao{db}
}

func (dao *SysRoleMenuDao) DeleteRoleMenuByRoleId(roleId int64) error {
	return dao.DB.Where("role_id = ?", roleId).Delete(&system.SysRoleMenu{}).Error
}

func (dao *SysRoleMenuDao) DeleteRoleMenu(ids []int64) error {
	return dao.DB.Where("role_id in (?)", ids).Delete(&system.SysRoleMenu{}).Error
}

func (dao *SysRoleMenuDao) BatchRoleMenu(list []*system.SysRoleMenu) error {
	return dao.DB.Create(&list).Error
}

func (dao *SysRoleMenuDao) CheckMenuExistRole(menuId int64) (count int64, err error) {
	err = dao.DB.Model(&system.SysRoleMenu{}).Where("menu_id = ?", menuId).Count(&count).Error
	return
}

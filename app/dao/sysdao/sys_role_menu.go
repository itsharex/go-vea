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

func (dao *SysRoleMenuDao) BatchRoleMenu(list []*system.SysRoleMenu) error {
	return dao.DB.Create(&list).Error
}

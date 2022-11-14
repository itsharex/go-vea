package sysdao

import (
	"context"
	"go-vea/app/model/system"
	"go-vea/configs"
	"gorm.io/gorm"
)

type SysRoleDeptDao struct {
	*gorm.DB
}

func NewSysRoleDeptDao(ctx context.Context) *SysRoleDeptDao {
	return &SysRoleDeptDao{configs.GetDB(ctx)}
}

func NewSysRoleDeptDaoByDB(db *gorm.DB) *SysRoleDeptDao {
	return &SysRoleDeptDao{db}
}

func (dao *SysRoleDeptDao) DeleteRoleDeptByRoleId(roleId int64) error {
	return dao.DB.Where("role_id = ?", roleId).Delete(&system.SysRoleDept{}).Error
}

func (dao *SysRoleDeptDao) BatchRoleDept(list []*system.SysRoleDept) error {
	return dao.DB.Create(&list).Error
}

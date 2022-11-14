package sysdao

import (
	"context"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/configs"
	"gorm.io/gorm"
)

type SysUserRoleDao struct {
	*gorm.DB
}

func NewSysUserRoleDao(ctx context.Context) *SysUserRoleDao {
	return &SysUserRoleDao{configs.GetDB(ctx)}
}

func NewSysUserRoleDaoByDB(db *gorm.DB) *SysUserRoleDao {
	return &SysUserRoleDao{db}
}

func (dao *SysUserRoleDao) SelectByUserId(userId int64) (roleIds []int64, err error) {
	err = dao.DB.Model(&system.SysUserRole{}).Select("role_id").Where("user_id = ?", userId).Find(&roleIds).Error
	return
}

func (dao *SysUserRoleDao) BatchUserRole(list []*system.SysUserRole) error {
	return dao.DB.Create(&list).Error
}

func (dao *SysUserRoleDao) DeleteUserRoleByUserId(userId int64) error {
	return dao.DB.Where("user_id = ?", userId).Delete(&system.SysUserRole{}).Error
}

func (dao *SysUserRoleDao) DeleteUserRole(ids []int64) error {
	return dao.DB.Where("user_id in (?)", ids).Delete(&system.SysUserRole{}).Error
}

func (dao *SysUserRoleDao) DeleteUserRoleInfo(ur *system.SysUserRole) error {
	return dao.DB.Where("user_id = ? and role_id = ? ", ur.UserID, ur.RoleID).Delete(&system.SysUserRole{}).Error
}

func (dao *SysUserRoleDao) BatchDeleteAuthUser(ur *request.SysUserRole) error {
	return dao.DB.Where("role_id = ? and user_id in (?)", ur.RoleID, ur.UserIds).Delete(&system.SysUserRole{}).Error
}

func (dao *SysUserRoleDao) BatchInsertAuthUser(urList []*system.SysUserRole) error {
	return dao.DB.Create(&urList).Error
}

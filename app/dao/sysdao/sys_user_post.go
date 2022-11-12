package sysdao

import (
	"context"
	"go-vea/app/model/system"
	"go-vea/configs"
	"gorm.io/gorm"
)

type SysUserPostDao struct {
	*gorm.DB
}

func NewSysUserPostDao(ctx context.Context) *SysUserPostDao {
	return &SysUserPostDao{configs.GetDB(ctx)}
}

func NewSysUserPostDaoByDB(db *gorm.DB) *SysUserPostDao {
	return &SysUserPostDao{db}
}

func (dao *SysUserPostDao) BatchUserPost(list []*system.SysUserPost) (err error) {
	return dao.DB.Create(&list).Error
}

func (dao *SysUserPostDao) DeleteUserPostByUserId(userId int64) error {
	return dao.DB.Where("user_id = ?", userId).Delete(&system.SysUserPost{}).Error
}

func (dao *SysUserPostDao) DeleteUserPost(ids []int64) error {
	return dao.DB.Where("user_id in (?)", ids).Delete(&system.SysUserPost{}).Error
}

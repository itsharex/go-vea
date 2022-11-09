package sysdao

import (
	"context"
	"go-web-template/app/common/e"
	"go-web-template/app/common/page"
	"go-web-template/app/model/system"
	"go-web-template/app/model/system/request"
	"go-web-template/configs"
	"gorm.io/gorm"
)

type SysUserDao struct {
	*gorm.DB
}

func NewSysUserDao(ctx context.Context) *SysUserDao {
	return &SysUserDao{configs.GetDB(ctx)}
}

func NewSysUserDaoByDB(db *gorm.DB) *SysUserDao {
	return &SysUserDao{db}
}

func (dao *SysUserDao) SelectList(sysUser *request.SysUser) (p *page.Pagination, err error) {
	var userLIst []*system.SysUser
	p = new(page.Pagination)

	if sysUser.DeptID != 0 {
		dao.DB = dao.DB.Where("dept_id = ?", sysUser.DeptID)
	}
	if sysUser.UserName != "" {
		dao.DB = dao.DB.Where("user_name = ?", sysUser.UserName)
	}
	if sysUser.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysUser.Status)
	}
	if sysUser.Phonenumber != "" {
		dao.DB = dao.DB.Where("phonenumber = ?", sysUser.Phonenumber)
	}

	if sysUser.OpenPage {
		p.PageNum = sysUser.PageNum
		p.PageSize = sysUser.PageSize
		err = dao.DB.Scopes(page.SelectPage(userLIst, p, dao.DB)).Find(&userLIst).Error
	} else {
		err = dao.DB.Find(&userLIst).Error
	}
	p.Rows = userLIst
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}

func (dao *SysUserDao) SelectAll(sysUser *request.SysUser) (list []system.SysUser, err error) {
	if sysUser.UserID != 0 {
		dao.DB = dao.DB.Where("user_id = ?", sysUser.UserID)
	}
	if sysUser.UserName != "" {
		dao.DB = dao.DB.Where("user_name = ?", sysUser.UserName)
	}
	if sysUser.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysUser.Status)
	}
	if sysUser.DataScope != "" {
		// todo
	}

	err = dao.DB.Where("del_flag = '0'").Find(&list).Error
	return
}

func (dao *SysUserDao) SelectUserByUserName(username string) (sysUser *system.SysUser, err error) {
	err = dao.DB.Model(&system.SysUser{}).Where("user_name=?", username).First(&sysUser).Error
	return
}

func (dao *SysUserDao) SelectById(id int64) (sysUser *system.SysUser, err error) {
	err = dao.DB.Model(&sysUser).Where("user_id=?", id).Find(&sysUser).Error
	return
}

func (dao *SysUserDao) Insert(sysUser *system.SysUser) error {
	return dao.DB.Model(&system.SysUser{}).Create(sysUser).Error
}

func (dao *SysUserDao) UpdateById(sysUser *system.SysUser) error {
	return dao.DB.Save(sysUser).Error
}

func (dao *SysUserDao) DeleteById(id int64) error {
	return dao.DB.Where("user_id = ?", id).Delete(&system.SysUser{}).Error
}

func (dao *SysUserDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("user_id in (?)", ids).Delete(&system.SysUser{}).Error
}

func (dao *SysUserDao) CheckUserNameUnique(roleName string) (count int64, err error) {
	err = dao.DB.Model(&system.SysUser{}).Where("user_name = ?", roleName).Where("del_flag = '0'").Count(&count).Error
	return
}

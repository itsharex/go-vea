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

type SysPostDao struct {
	*gorm.DB
}

func NewSysPostDao(ctx context.Context) *SysPostDao {
	return &SysPostDao{configs.GetDB(ctx)}
}

func NewSysPostDaoByDB(db *gorm.DB) *SysPostDao {
	return &SysPostDao{db}
}

func (dao *SysPostDao) SelectList(sysPost *request.SysPost) (p *page.Pagination, err error) {
	var postList []*system.SysPost
	p = new(page.Pagination)

	if sysPost.PostCode != "" {
		dao.DB = dao.DB.Where("post_code = ?", sysPost.PostCode)
	}
	if sysPost.PostName != "" {
		dao.DB = dao.DB.Where("post_name = ?", sysPost.PostName)
	}
	if sysPost.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysPost.Status)
	}

	if sysPost.OpenPage {
		p.PageNum = sysPost.PageNum
		p.PageSize = sysPost.PageSize
		err = dao.DB.Scopes(page.SelectPage(postList, p, dao.DB)).Find(&postList).Error
	} else {
		err = dao.DB.Find(&postList).Error
	}
	p.Rows = postList
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}

func (dao *SysPostDao) SelectAll(sysPost *request.SysPost) (list []system.SysPost, err error) {
	if sysPost.PostCode != "" {
		dao.DB = dao.DB.Where("post_code = ?", sysPost.PostCode)
	}
	if sysPost.PostName != "" {
		dao.DB = dao.DB.Where("post_name = ?", sysPost.PostName)
	}
	if sysPost.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysPost.Status)
	}

	return
}

func (dao *SysPostDao) SelectById(id int64) (sysPost *system.SysPost, err error) {
	err = dao.DB.Model(&sysPost).Where("post_id=?", id).Find(&sysPost).Error
	return
}

func (dao *SysPostDao) Insert(sysPost *system.SysPost) error {
	return dao.DB.Model(&system.SysPost{}).Create(sysPost).Error
}

func (dao *SysPostDao) UpdateById(sysPost *system.SysPost) error {
	return dao.DB.Save(sysPost).Error
}

func (dao *SysPostDao) DeleteById(id int64) error {
	return dao.DB.Where("post_id = ?", id).Delete(&system.SysPost{}).Error
}

func (dao *SysPostDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("post_id in (?)", ids).Delete(&system.SysPost{}).Error
}

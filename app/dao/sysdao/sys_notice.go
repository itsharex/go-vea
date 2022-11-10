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

type SysNoticeDao struct {
	*gorm.DB
}

func NewSysNoticeDao(ctx context.Context) *SysNoticeDao {
	return &SysNoticeDao{configs.GetDB(ctx)}
}

func NewSysNoticeDaoByDB(db *gorm.DB) *SysNoticeDao {
	return &SysNoticeDao{db}
}

func (dao *SysNoticeDao) SelectList(sysNotice *request.SysNotice) (p *page.Pagination, err error) {
	var noticeList []*system.SysNotice
	p = new(page.Pagination)

	if sysNotice.NoticeTitle != "" {
		dao.DB = dao.DB.Where("notice_title = ?", sysNotice.NoticeTitle)
	}
	if sysNotice.NoticeType != "" {
		dao.DB = dao.DB.Where("notice_type = ?", sysNotice.NoticeType)
	}

	if sysNotice.OpenPage {
		p.PageNum = sysNotice.PageNum
		p.PageSize = sysNotice.PageSize
		err = dao.DB.Scopes(page.SelectPage(noticeList, p, dao.DB)).Find(&noticeList).Error
	} else {
		err = dao.DB.Find(&noticeList).Error
	}
	p.Rows = noticeList
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}

func (dao *SysNoticeDao) SelectAll(sysNotice *request.SysNotice) (list []system.SysNotice, err error) {
	if sysNotice.NoticeTitle != "" {
		dao.DB = dao.DB.Where("notice_title = ?", sysNotice.NoticeTitle)
	}
	if sysNotice.NoticeType != "" {
		dao.DB = dao.DB.Where("notice_type = ?", sysNotice.NoticeType)
	}

	return
}

func (dao *SysNoticeDao) SelectById(id int64) (sysNotice *system.SysNotice, err error) {
	err = dao.DB.Model(&sysNotice).Where("notice_id=?", id).Find(&sysNotice).Error
	return
}

func (dao *SysNoticeDao) Insert(sysNotice *system.SysNotice) error {
	return dao.DB.Model(&system.SysNotice{}).Create(sysNotice).Error
}

func (dao *SysNoticeDao) UpdateById(sysNotice *system.SysNotice) error {
	return dao.DB.Save(sysNotice).Error
}

func (dao *SysNoticeDao) DeleteById(id int64) error {
	return dao.DB.Where("notice_id = ?", id).Delete(&system.SysNotice{}).Error
}

func (dao *SysNoticeDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("notice_id in (?)", ids).Delete(&system.SysNotice{}).Error
}

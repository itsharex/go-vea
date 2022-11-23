package monitordao

import (
	"context"
	"go-vea/app/common/e"
	"go-vea/app/common/page"
	"go-vea/app/model/monitor"
	"go-vea/app/model/monitor/request"
	"go-vea/configs"
	"gorm.io/gorm"
)

type SysOperLogDao struct {
	*gorm.DB
}

func NewSysOperLogDao(ctx context.Context) *SysOperLogDao {
	return &SysOperLogDao{configs.GetDB(ctx)}
}

func NewSysOperLogDaoByDB(db *gorm.DB) *SysOperLogDao {
	return &SysOperLogDao{db}
}

func (dao *SysOperLogDao) SelectList(operLog *request.SysOperLog) (p *page.Pagination, err error) {
	var configList []*monitor.SysOperLog
	p = new(page.Pagination)

	if operLog.Title != "" {
		dao.DB = dao.DB.Where("title = ?", operLog.Title)
	}
	if operLog.BusinessType != 0 {
		dao.DB = dao.DB.Where("business_type = ?", operLog.BusinessType)
	}
	if operLog.BusinessTypes != nil {
		dao.DB = dao.DB.Where("business_type in (?)", operLog.BusinessTypes)
	}
	if operLog.Status != 0 {
		dao.DB = dao.DB.Where("status = ?", operLog.Status)
	}
	if operLog.OperName != "" {
		dao.DB = dao.DB.Where("oper_name = ?", operLog.OperName)
	}

	if operLog.OpenPage {
		p.PageNum = operLog.PageNum
		p.PageSize = operLog.PageSize
		err = dao.DB.Scopes(page.SelectPage(configList, p, dao.DB)).Find(&configList).Error
	} else {
		err = dao.DB.Find(&configList).Error
	}
	p.Rows = configList
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}

func (dao *SysOperLogDao) SelectAll() (list []monitor.SysOperLog, err error) {
	err = dao.DB.Find(&list).Error
	return
}

func (dao *SysOperLogDao) SelectById(id int64) (operLog *monitor.SysOperLog, err error) {
	err = dao.DB.Model(&operLog).Where("oper_id=?", id).Find(&operLog).Error
	return
}

func (dao *SysOperLogDao) Insert(operLog *monitor.SysOperLog) error {
	return dao.DB.Model(&monitor.SysOperLog{}).Create(operLog).Error
}

func (dao *SysOperLogDao) UpdateById(operLog *monitor.SysOperLog) error {
	return dao.DB.Updates(operLog).Error
}

func (dao *SysOperLogDao) DeleteById(id int64) error {
	return dao.DB.Where("oper_id = ?", id).Delete(&monitor.SysOperLog{}).Error
}

func (dao *SysOperLogDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("oper_id in (?)", ids).Delete(&monitor.SysOperLog{}).Error
}

func (dao *SysOperLogDao) CleanOperLog() error {
	return dao.DB.Exec("truncate table sys_oper_log").Error
}

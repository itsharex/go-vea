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

type SysLoginLogDao struct {
	*gorm.DB
}

func NewSysLoginLogDao(ctx context.Context) *SysLoginLogDao {
	return &SysLoginLogDao{configs.GetDB(ctx)}
}

func NewSysLoginLogDaoByDB(db *gorm.DB) *SysLoginLogDao {
	return &SysLoginLogDao{db}
}

func (dao *SysLoginLogDao) SelectList(sysLoginLog *request.SysLoginLog) (p *page.Pagination, err error) {
	var configList []*monitor.SysLoginLog
	p = new(page.Pagination)

	if sysLoginLog.Ipaddr != "" {
		dao.DB = dao.DB.Where("ipaddr = ?", sysLoginLog.Ipaddr)
	}
	if sysLoginLog.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysLoginLog.Status)
	}
	if sysLoginLog.Username != "" {
		dao.DB = dao.DB.Where("username = ?", sysLoginLog.Username)
	}

	if sysLoginLog.OpenPage {
		p.PageNum = sysLoginLog.PageNum
		p.PageSize = sysLoginLog.PageSize
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

func (dao *SysLoginLogDao) SelectAll() (list []monitor.SysLoginLog, err error) {
	err = dao.DB.Find(&list).Error
	return
}

func (dao *SysLoginLogDao) SelectById(id int64) (loginLog *monitor.SysLoginLog, err error) {
	err = dao.DB.Model(&loginLog).Where("id=?", id).Find(&loginLog).Error
	return
}

func (dao *SysLoginLogDao) Insert(loginLog *monitor.SysLoginLog) error {
	return dao.DB.Model(&monitor.SysLoginLog{}).Create(loginLog).Error
}

func (dao *SysLoginLogDao) UpdateById(loginLog *monitor.SysLoginLog) error {
	return dao.DB.Updates(loginLog).Error
}

func (dao *SysLoginLogDao) DeleteById(id int64) error {
	return dao.DB.Where("id = ?", id).Delete(&monitor.SysLoginLog{}).Error
}

func (dao *SysLoginLogDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("id in (?)", ids).Delete(&monitor.SysLoginLog{}).Error
}

func (dao *SysLoginLogDao) CleanLoginLog() error {
	return dao.DB.Exec("truncate table sys_login_log").Error
}

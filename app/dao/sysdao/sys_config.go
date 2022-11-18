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

type SysConfigDao struct {
	*gorm.DB
}

func NewSysConfigDao(ctx context.Context) *SysConfigDao {
	return &SysConfigDao{configs.GetDB(ctx)}
}

func NewSysConfigDaoByDB(db *gorm.DB) *SysConfigDao {
	return &SysConfigDao{db}
}

func (dao *SysConfigDao) SelectList(sysConfig *request.SysConfig) (p *page.Pagination, err error) {
	var configList []*system.SysConfig
	p = new(page.Pagination)

	if sysConfig.ConfigName != "" {
		dao.DB = dao.DB.Where("config_name = ?", sysConfig.ConfigName)
	}
	if sysConfig.ConfigType != "" {
		dao.DB = dao.DB.Where("config_type = ?", sysConfig.ConfigType)
	}
	if sysConfig.ConfigKey != "" {
		dao.DB = dao.DB.Where("config_key = ?", sysConfig.ConfigKey)
	}

	if sysConfig.OpenPage {
		p.PageNum = sysConfig.PageNum
		p.PageSize = sysConfig.PageSize
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

func (dao *SysConfigDao) SelectAll() (list []system.SysConfig, err error) {
	err = dao.DB.Find(&list).Error
	return
}

func (dao *SysConfigDao) SelectById(id int64) (config *system.SysConfig, err error) {
	err = dao.DB.Model(&config).Where("config_id=?", id).Find(&config).Error
	return
}

func (dao *SysConfigDao) Insert(config *system.SysConfig) error {
	return dao.DB.Model(&system.SysConfig{}).Create(config).Error
}

func (dao *SysConfigDao) UpdateById(config *system.SysConfig) error {
	return dao.DB.Updates(config).Error
}

func (dao *SysConfigDao) DeleteById(id int64) error {
	return dao.DB.Where("config_id = ?", id).Delete(&system.SysConfig{}).Error
}

func (dao *SysConfigDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("config_id in (?)", ids).Delete(&system.SysConfig{}).Error
}

func (dao *SysConfigDao) CheckConfigKeyUnique(configKey string) (sysConfig *system.SysConfig, err error) {
	err = dao.DB.Model(&system.SysConfig{}).Where("config_key=?", configKey).First(&sysConfig).Error
	return
}

func (dao *SysConfigDao) SelectSysConfigByKey(configKey string) (config *system.SysConfig, err error) {
	err = dao.DB.Model(&config).Where("config_key=?", configKey).First(&config).Error
	return
}

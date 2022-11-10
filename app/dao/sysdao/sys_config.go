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
	// Save 保存更新数据库中的值。如果值不包含一个匹配的主键，值将被插入。
	return dao.DB.Save(config).Error
}

func (dao *SysConfigDao) DeleteById(id int64) error {
	return dao.DB.Where("config_id = ?", id).Delete(&system.SysConfig{}).Error
}

func (dao *SysConfigDao) DeleteByIds(ids []int64) error {
	// Delete 删除符合给定条件的值。如果值包含主键，它就会被包括在条件中。
	// 如果值包含一个delete_at字段，那么Delete将执行一个软删除，如果空的话，将delete_at设置为当前时间。
	//return dao.DB.Delete(&system.SysConfig{}, "config_id in (?)", ids).Error
	return dao.DB.Where("config_id in (?)", ids).Delete(&system.SysConfig{}).Error
}

func (dao *SysConfigDao) CheckConfigKeyUnique(configKey string) (count int64, err error) {
	err = dao.DB.Model(&system.SysConfig{}).Where("config_key=?", configKey).Count(&count).Error
	return
}

func (dao *SysConfigDao) SelectSysConfigByKey(configKey string) (config *system.SysConfig, err error) {
	err = dao.DB.Model(&config).Where("config_key=?", configKey).First(&config).Error
	return
}

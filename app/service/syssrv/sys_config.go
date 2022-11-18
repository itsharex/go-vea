package syssrv

import (
	"context"
	"errors"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
	"gorm.io/gorm"
)

type SysConfigService struct{}

var SysConfigSrv = new(SysConfigService)

func (*SysConfigService) SelectSysConfigList(ctx context.Context, sysConfig *request.SysConfig) (*page.Pagination, error) {
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	data, err := sysConfigDao.SelectList(sysConfig)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysConfigService) SelectSysConfigById(ctx context.Context, configId int64) (*system.SysConfig, error) {
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	data, err := sysConfigDao.SelectById(configId)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (*SysConfigService) SelectSysConfigByKey(ctx context.Context, configKey string) (configValue string, err error) {
	configValue, err = global.Redis.Get(ctx, getConfigKey(configKey)).Result()
	if configValue != "" && err == nil {
		return configValue, nil
	}
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	data, e := sysConfigDao.SelectSysConfigByKey(configKey)
	if data != nil && err == nil {
		global.Redis.Set(ctx, getConfigKey(data.ConfigKey), data.ConfigValue, 0)
		return data.ConfigValue, e
	}
	return "", nil
}

func (s *SysConfigService) SelectCaptchaEnabled(ctx context.Context) (bool, error) {
	value, err := s.SelectSysConfigByKey(ctx, "sys.account.captchaEnabled")
	if err != nil {
		return false, err
	}
	return value == "true", nil
}

func (s *SysConfigService) AddSysConfig(ctx context.Context, sysConfig *system.SysConfig) error {
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	configKeyUnique := s.checkConfigKeyUnique(ctx, sysConfig)
	if !configKeyUnique {
		global.Logger.Error("新增失败！已存在该配置key")
		return errors.New("新增失败！已存在该配置key")
	}
	err := sysConfigDao.Insert(sysConfig)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	global.Redis.Set(ctx, getConfigKey(sysConfig.ConfigKey), sysConfig.ConfigValue, 0)
	return nil
}

func (s *SysConfigService) UpdateSysConfig(ctx context.Context, sysConfig *system.SysConfig) error {
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	configKeyUnique := s.checkConfigKeyUnique(ctx, sysConfig)
	if !configKeyUnique {
		global.Logger.Error("修改失败！已存在该配置key")
		return errors.New("修改失败！已存在该配置key")
	}
	err := sysConfigDao.UpdateById(sysConfig)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	global.Redis.Set(ctx, getConfigKey(sysConfig.ConfigKey), sysConfig.ConfigValue, 0)
	return nil
}

func (*SysConfigService) DeleteSysConfigByIds(ctx context.Context, ids []int64) error {
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	for _, id := range ids {
		data, _ := sysConfigDao.SelectById(id)
		if data.ConfigType == "Y" {
			return errors.New("系统内置，无法删除")
		}
		err := sysConfigDao.DeleteById(id)
		global.Redis.Del(ctx, getConfigKey(data.ConfigKey))
		if err != nil {
			global.Logger.Error(err)
			return err
		}
	}
	return nil
}

func (*SysConfigService) ResetConfigCache(ctx context.Context) {
	clearConfigCache(ctx)
	loadingConfigCache(ctx)
}

func (*SysConfigService) checkConfigKeyUnique(ctx context.Context, sysConfig *system.SysConfig) bool {
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	data, err := sysConfigDao.CheckConfigKeyUnique(sysConfig.ConfigKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.ConfigID != 0 && data.ConfigID != sysConfig.ConfigID {
		return false
	}
	return true
}

func loadingConfigCache(ctx context.Context) {
	sysConfigDao := sysdao.NewSysConfigDao(ctx)
	list, _ := sysConfigDao.SelectAll()
	for _, sysConfig := range list {
		global.Redis.Set(ctx, getConfigKey(sysConfig.ConfigKey), sysConfig.ConfigValue, 0)
	}
}

func clearConfigCache(ctx context.Context) {
	keys, _ := global.Redis.Keys(ctx, getConfigKey("*")).Result()
	global.Redis.Del(ctx, keys...)
}

func getConfigKey(configKey string) string {
	return "sys_config:" + configKey
}

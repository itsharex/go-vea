package syssrv

import (
	"context"
	"encoding/json"
	"errors"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
	"gorm.io/gorm"
)

type SysDictTypeService struct{}

var SysDictTypeSrv = new(SysDictTypeService)

func (*SysDictTypeService) InitDictCache(ctx context.Context) {
	err := loadingDictCache(ctx)
	if err != nil {
		global.Logger.Error("初始化字典缓存失败")
	}
}

func (*SysDictTypeService) SelectDictTypeList(ctx context.Context, dictType request.SysDictType) (*page.Pagination, error) {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	data, err := sysDictTypeDao.SelectList(dictType)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysDictTypeService) SelectDictTypeAll(ctx context.Context) ([]system.SysDictType, error) {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	data, err := sysDictTypeDao.SelectAll()
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return data, err
}

func (*SysDictTypeService) SelectDictTypeById(ctx context.Context, dictId int64) (*system.SysDictType, error) {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	data, err := sysDictTypeDao.SelectById(dictId)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return data, err
}

func (*SysDictTypeService) SelectDictTypeByType(ctx context.Context, dictType request.SysDictType) (*page.Pagination, error) {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	data, err := sysDictTypeDao.SelectList(dictType)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return data, err
}

func (s *SysDictTypeService) DeleteDictTypeByIds(ctx context.Context, ids []int64) error {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	for _, dictId := range ids {
		sysDictType, _ := sysDictTypeDao.SelectById(dictId)
		used := s.countDictDataByType(ctx, sysDictType.DictType)
		if sysDictType != nil && used {
			return errors.New("所选字典已分配,不能删除")
		}
	}
	err := sysDictTypeDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (s *SysDictTypeService) AddDictType(ctx context.Context, sysDictType *system.SysDictType) error {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	dictTypeUnique := s.checkDictTypeUnique(ctx, sysDictType)
	if !dictTypeUnique {
		global.Logger.Error("新增失败！已存在该字典类型")
		return errors.New("新增失败！已存在该字典类型")
	}
	err := sysDictTypeDao.Insert(sysDictType)
	// todo 更新缓存
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (s *SysDictTypeService) UpdateDictType(ctx context.Context, sysDictType *system.SysDictType) error {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	dictTypeUnique := s.checkDictTypeUnique(ctx, sysDictType)
	if !dictTypeUnique {
		global.Logger.Error("修改失败！已存在该字典类型")
		return errors.New("修改失败！已存在该字典类型")
	}
	err := sysDictTypeDao.UpdateById(sysDictType)
	// todo 更新缓存
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysDictTypeService) checkDictTypeUnique(ctx context.Context, sysDictType *system.SysDictType) bool {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	data, err := sysDictTypeDao.CheckDictTypeUnique(sysDictType.DictType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.DictID != 0 && data.DictID != sysDictType.DictID {
		return false
	}
	return true
}

func (*SysDictTypeService) countDictDataByType(ctx context.Context, dictType string) bool {
	sysDictTypeDao := sysdao.NewSysDictTypeDao(ctx)
	count, err := sysDictTypeDao.CountDictDataByType(dictType)
	if err != nil {
		return true
	}
	return count > 0
}

func (*SysDictTypeService) ResetDictCache(ctx context.Context) (err error) {
	err = clearDictCache(ctx)
	err = loadingDictCache(ctx)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func loadingDictCache(ctx context.Context) (err error) {
	sysDictDataDao := sysdao.NewSysDictDataDao(ctx)
	list, err := sysDictDataDao.SelectAll()
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	dictMap := make(map[string][]system.SysDictData)
	for _, dictData := range list {
		dictMap[dictData.DictType] = append(dictMap[dictData.DictType], dictData)
	}
	for key, data := range dictMap {
		jsonData, _ := json.Marshal(data)
		err = global.Redis.Set(ctx, getDictKey(key), jsonData, 0).Err()
	}
	return nil
}

func clearDictCache(ctx context.Context) error {
	keys, _ := global.Redis.Keys(ctx, getDictKey("*")).Result()
	err := global.Redis.Del(ctx, keys...).Err()
	return err
}

func getDictKey(dictKey string) string {
	return "sys_dict:" + dictKey
}

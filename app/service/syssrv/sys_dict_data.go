package syssrv

import (
	"context"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
)

type SysDictDataService struct{}

var SysDictDataSrv = new(SysDictDataService)

// SelectDictDataList 根据条件分页查询字典数据
func (*SysDictDataService) SelectDictDataList(ctx context.Context, dictData request.SysDictData) (*page.Pagination, error) {
	sysDictDataDao := sysdao.NewSysDictDataDao(ctx)
	data, err := sysDictDataDao.SelectList(dictData)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

// SelectDictDataById 根据字典数据ID查询信息
func (*SysDictDataService) SelectDictDataById(ctx context.Context, dictCode int64) (*system.SysDictData, error) {
	sysDictDataDao := sysdao.NewSysDictDataDao(ctx)
	data, err := sysDictDataDao.SelectById(dictCode)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return data, err
}

// AddDictData 新增
func (*SysDictDataService) AddDictData(ctx context.Context, dictData *system.SysDictData) error {
	sysDictDataDao := sysdao.NewSysDictDataDao(ctx)
	err := sysDictDataDao.Insert(dictData)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

// UpdateDictData 更新
func (s *SysDictDataService) UpdateDictData(ctx context.Context, dictData *system.SysDictData) error {
	sysDictDataDao := sysdao.NewSysDictDataDao(ctx)
	err := sysDictDataDao.UpdateById(dictData)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

// DeleteDictDataByIds 单个删除/批量删除
func (s *SysDictDataService) DeleteDictDataByIds(ctx context.Context, ids []int64) error {
	sysDictDataDao := sysdao.NewSysDictDataDao(ctx)
	err := sysDictDataDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

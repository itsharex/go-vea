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

type SysDictDataDao struct {
	*gorm.DB
}

func NewSysDictDataDao(ctx context.Context) *SysDictDataDao {
	return &SysDictDataDao{configs.GetDB(ctx)}
}

func NewSysDictDataDaoByDB(db *gorm.DB) *SysDictDataDao {
	return &SysDictDataDao{db}
}

func (dao *SysDictDataDao) SelectList(dictData request.SysDictData) (p *page.Pagination, err error) {
	var dictDataList []*system.SysDictData
	p = new(page.Pagination)

	if dictData.DictType != "" {
		dao.DB = dao.DB.Where("dict_type = ?", dictData.DictType)
	}
	if dictData.DictLabel != "" {
		dao.DB = dao.DB.Where("dict_label = ?", dictData.DictLabel)
	}
	if dictData.Status != "" {
		dao.DB = dao.DB.Where("status = ?", dictData.Status)
	}
	dao.DB.Order("dict_sort asc")

	if dictData.OpenPage {
		p.PageNum = dictData.PageNum
		p.PageSize = dictData.PageSize
		err = dao.DB.Scopes(page.SelectPage(dictDataList, p, dao.DB)).Find(&dictDataList).Error
	} else {
		err = dao.DB.Find(&dictDataList).Error
	}
	p.Rows = dictDataList
	if err != nil {
		p.Code = e.ERROR
		p.Msg = err.Error()
		return p, err
	}
	return p, err
}

func (dao *SysDictDataDao) SelectAll() (list []system.SysDictData, err error) {
	err = dao.DB.Where("status = 0").Order("dict_sort asc").Find(&list).Error
	return
}

func (dao *SysDictDataDao) SelectById(dictCode int64) (dictData *system.SysDictData, err error) {
	err = dao.DB.Where("dict_code = ?", dictCode).Find(&dictData).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *SysDictDataDao) Insert(dictData *system.SysDictData) error {
	return dao.DB.Create(dictData).Error
}

func (dao *SysDictDataDao) UpdateById(dictData *system.SysDictData) error {
	return dao.DB.Save(dictData).Error
}

func (dao *SysDictDataDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("dict_code in (?)", ids).Delete(&system.SysDictData{}).Error
}

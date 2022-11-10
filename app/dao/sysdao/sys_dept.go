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

type SysDeptDao struct {
	*gorm.DB
}

func NewSysDeptDao(ctx context.Context) *SysDeptDao {
	return &SysDeptDao{configs.GetDB(ctx)}
}

func NewSysDeptDaoByDB(db *gorm.DB) *SysDeptDao {
	return &SysDeptDao{db}
}

func (dao *SysDeptDao) SelectList(sysDept *request.SysDept) (p *page.Pagination, err error) {
	var postList []*system.SysDept
	p = new(page.Pagination)

	if sysDept.DeptID != 0 {
		dao.DB = dao.DB.Where("dept_id = ?", sysDept.DeptID)
	}
	if sysDept.ParentID != 0 {
		dao.DB = dao.DB.Where("parent_id = ?", sysDept.ParentID)
	}
	if sysDept.DeptName != "" {
		dao.DB = dao.DB.Where("dept_name = ?", sysDept.DeptName)
	}
	if sysDept.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysDept.Status)
	}

	dao.DB.Where("del_flag = '0'")

	if sysDept.OpenPage {
		p.PageNum = sysDept.PageNum
		p.PageSize = sysDept.PageSize
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

func (dao *SysDeptDao) SelectAll(sysDept *request.SysDept) (list []system.SysDept, err error) {
	if sysDept.DeptID != 0 {
		dao.DB = dao.DB.Where("dept_id = ?", sysDept.DeptID)
	}
	if sysDept.ParentID != 0 {
		dao.DB = dao.DB.Where("parent_id = ?", sysDept.ParentID)
	}
	if sysDept.DeptName != "" {
		dao.DB = dao.DB.Where("dept_name = ?", sysDept.DeptName)
	}
	if sysDept.Status != "" {
		dao.DB = dao.DB.Where("status = ?", sysDept.Status)
	}

	dao.DB.Where("del_flag = '0'")

	return
}

func (dao *SysDeptDao) SelectById(id int64) (sysDept *system.SysDept, err error) {
	err = dao.DB.Model(&sysDept).Where("dept_id=?", id).Find(&sysDept).Error
	return
}

func (dao *SysDeptDao) Insert(sysDept *system.SysDept) error {
	return dao.DB.Model(&system.SysDept{}).Create(sysDept).Error
}

func (dao *SysDeptDao) UpdateById(sysDept *system.SysDept) error {
	return dao.DB.Save(sysDept).Error
}

func (dao *SysDeptDao) DeleteById(id int64) error {
	return dao.DB.Where("dept_id = ?", id).Delete(&system.SysDept{}).Error
}

func (dao *SysDeptDao) DeleteByIds(ids []int64) error {
	return dao.DB.Where("dept_id in (?)", ids).Delete(&system.SysDept{}).Error
}

package syssrv

import (
	"context"
	"errors"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
)

type SysDeptService struct{}

var SysDeptSrv = new(SysDeptService)

func (*SysDeptService) GetSysDeptList(ctx context.Context, sysDept *request.SysDept) (*page.Pagination, error) {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	data, err := sysDeptDao.SelectList(sysDept)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return data, err
}

func (*SysDeptService) GetSysDeptById(ctx context.Context, roleId int64) (*system.SysDept, error) {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	data, err := sysDeptDao.SelectById(roleId)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	return data, err
}

func (*SysDeptService) AddSysDept(ctx context.Context, sysDept *system.SysDept) error {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	err := sysDeptDao.Insert(sysDept)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysDeptService) UpdateDeptById(ctx context.Context, sysDept *system.SysDept) error {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	err := sysDeptDao.UpdateById(sysDept)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysDeptService) DeleteSysDeptByIds(ctx context.Context, ids []int64) error {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	err := sysDeptDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysDeptService) GetDeptTreeList(ctx context.Context, sysDept *request.SysDept) ([]*system.SysDept, error) {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	data, err := sysDeptDao.SelectList(sysDept)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	// 类型转换
	depts, ok := data.Rows.([]*system.SysDept)
	if !ok {
		global.Logger.Error("类型转换错误")
		return nil, errors.New("类型转换错误")
	}
	treeList := buildDeptTree(depts)
	return treeList, err
}

func buildDeptTree(depts []*system.SysDept) []*system.SysDept {
	deptMap := make(map[int64]*system.SysDept)
	for i, dept := range depts {
		dept.ArrIdx = i
		deptMap[dept.DeptID] = dept
	}

	var resList []*system.SysDept
	var childList []*system.SysDept

	for _, dept := range depts {

		parent, ok := deptMap[dept.ParentID]

		if ok {
			if len(parent.Children) == 0 {
				childList = depts[parent.ArrIdx].Children
				if childList == nil {
					childList = []*system.SysDept{}
				}
				childList = append(childList, dept)
				depts[parent.ArrIdx].Children = childList
			} else {
				depts[parent.ArrIdx].Children = append(depts[parent.ArrIdx].Children, dept)
			}
		}

		if dept.ParentID == 0 {
			resList = append(resList, dept)
		}
	}
	return resList
}

func hasChildByDeptId(ctx context.Context) bool {
	return false
}

func checkDeptNameUnique(ctx context.Context) bool {
	return false
}

func checkDeptExistUser(ctx context.Context) bool {
	return false
}

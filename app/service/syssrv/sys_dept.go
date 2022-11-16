package syssrv

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
	"go-vea/util"
	"gorm.io/gorm"
	"strconv"
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
	unique := checkDeptNameUnique(ctx, sysDept)
	if !unique {
		global.Logger.Error("新增部门'" + sysDept.DeptName + "'失败，部门名称已存在")
		return errors.New("新增部门'" + sysDept.DeptName + "'失败，部门名称已存在")
	}
	err := sysDeptDao.Insert(sysDept)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysDeptService) UpdateDeptById(ctx context.Context, sysDept *system.SysDept) error {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	unique := checkDeptNameUnique(ctx, sysDept)
	normalStatus := selectNormalChildrenDeptById(ctx, sysDept.DeptID)
	if !unique {
		global.Logger.Error("修改部门'" + sysDept.DeptName + "'失败，部门名称已存在")
		return errors.New("修改部门'" + sysDept.DeptName + "'失败，部门名称已存在")
	} else if sysDept.ParentID == sysDept.DeptID {
		global.Logger.Error("修改部门'" + sysDept.DeptName + "'失败，上级部门不能是自己")
		return errors.New("修改部门'" + sysDept.DeptName + "'失败，上级部门不能是自己")
	} else if sysDept.Status == "1" && normalStatus {
		global.Logger.Error("该部门包含未停用的子部门")
		return errors.New("该部门包含未停用的子部门")
	}
	err := sysDeptDao.UpdateById(sysDept)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysDeptService) DeleteSysDeptById(ctx context.Context, deptId int64) error {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	hasChild := hasChildByDeptId(ctx, deptId)
	existUser := checkDeptExistUser(ctx, deptId)
	if hasChild {
		global.Logger.Error("存在下级部门,不允许删除")
		return errors.New("存在下级部门,不允许删除")
	}
	if existUser {
		global.Logger.Error("部门存在用户,不允许删除")
		return errors.New("部门存在用户,不允许删除")
	}
	err := sysDeptDao.DeleteById(deptId)
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
	deptList, ok := data.Rows.([]*system.SysDept)
	if !ok {
		global.Logger.Error("类型转换错误")
		return nil, errors.New("类型转换错误")
	}
	treeList := buildDeptTree(deptList)
	return treeList, err
}

func (*SysDeptService) GetDeptTreeListExcludeChild(ctx context.Context, sysDept *request.SysDept) ([]*system.SysDept, error) {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	data, err := sysDeptDao.SelectList(sysDept)
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	// 类型转换
	deptList, ok := data.Rows.([]*system.SysDept)
	if !ok {
		global.Logger.Error("类型转换错误")
		return nil, errors.New("类型转换错误")
	}
	// 排除子节点
	var res []*system.SysDept
	if sysDept.ExcludeDeptId != -1 {
		for _, dept := range deptList {
			c := util.Contains(dept.Ancestors, strconv.FormatInt(sysDept.ExcludeDeptId, 10))
			if !(dept.DeptID == sysDept.ExcludeDeptId || c) {
				res = append(res, dept)
			}
		}
	} else {
		res = deptList
	}
	// 构建树结构
	treeList := buildDeptTree(res)
	return treeList, err
}

func (*SysDeptService) SelectDeptListByRoleId(ctx *gin.Context, roleId int64) (deptIds []int64, err error) {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	sysRoleDao := sysdao.NewSysRoleDao(ctx)
	sysRole, err := sysRoleDao.SelectById(roleId)
	deptIds, err = sysDeptDao.SelectDeptListByRoleId(roleId, sysRole.IsDeptCheckStrictly(sysRole.DeptCheckStrictly))
	if err != nil {
		return nil, err
	}
	return deptIds, nil
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

func hasChildByDeptId(ctx context.Context, deptId int64) bool {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	count, err := sysDeptDao.HasChildByDeptId(deptId)
	if err != nil {
		return false
	}
	return count > 0
}

func checkDeptNameUnique(ctx context.Context, dept *system.SysDept) bool {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	data, err := sysDeptDao.CheckDeptNameUnique(dept)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.DeptID != 0 && data.DeptID != dept.DeptID {
		return false
	}
	return true
}

func checkDeptExistUser(ctx context.Context, deptId int64) bool {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	count, err := sysDeptDao.CheckDeptExistUser(deptId)
	if err != nil {
		return false
	}
	return count > 0
}

func selectNormalChildrenDeptById(ctx context.Context, deptId int64) bool {
	sysDeptDao := sysdao.NewSysDeptDao(ctx)
	count, err := sysDeptDao.SelectNormalChildrenDeptById(deptId)
	if err != nil {
		return false
	}
	return count > 0
}

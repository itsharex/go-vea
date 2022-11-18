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

type SysPostService struct{}

var SysPostSrv = new(SysPostService)

func (*SysPostService) GetSysPostList(ctx context.Context, sysPost *request.SysPost) (*page.Pagination, error) {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	data, err := sysPostDao.SelectList(sysPost)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysPostService) GetSysPostById(ctx context.Context, roleId int64) (*system.SysPost, error) {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	data, err := sysPostDao.SelectById(roleId)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysPostService) AddSysPost(ctx context.Context, sysPost *system.SysPost) error {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	postNameUnique := checkPostNameUnique(ctx, sysPost)
	postCodeUnique := checkPostCodeUnique(ctx, sysPost)
	if !postNameUnique {
		global.Logger.Error("新增岗位'" + sysPost.PostName + "'失败，岗位名称已存在")
		return errors.New("新增岗位'" + sysPost.PostName + "'失败，岗位名称已存在")
	}
	if !postCodeUnique {
		global.Logger.Error("新增岗位'" + sysPost.PostName + "'失败，岗位编码已存在")
		return errors.New("新增岗位'" + sysPost.PostName + "'失败，岗位编码已存在")
	}
	err := sysPostDao.Insert(sysPost)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysPostService) UpdatePostById(ctx context.Context, sysPost *system.SysPost) error {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	postNameUnique := checkPostNameUnique(ctx, sysPost)
	postCodeUnique := checkPostCodeUnique(ctx, sysPost)
	if !postNameUnique {
		global.Logger.Error("修改岗位'" + sysPost.PostName + "'失败，岗位名称已存在")
		return errors.New("修改岗位'" + sysPost.PostName + "'失败，岗位名称已存在")
	}
	if !postCodeUnique {
		global.Logger.Error("修改岗位'" + sysPost.PostName + "'失败，岗位编码已存在")
		return errors.New("修改岗位'" + sysPost.PostName + "'失败，岗位编码已存在")
	}
	err := sysPostDao.UpdateById(sysPost)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysPostService) DeleteSysPostByIds(ctx context.Context, ids []int64) error {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	err := sysPostDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysPostService) SelectPostAll(ctx context.Context) ([]*system.SysPost, error) {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	data, err := sysPostDao.SelectAll(&request.SysPost{})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (*SysPostService) SelectPostListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	data, err := sysPostDao.SelectPostListByUserId(userId)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func checkPostNameUnique(ctx context.Context, post *system.SysPost) bool {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	data, err := sysPostDao.CheckPostNameUnique(post.PostName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.PostID != 0 && data.PostID != post.PostID {
		return false
	}
	return true
}

func checkPostCodeUnique(ctx context.Context, post *system.SysPost) bool {
	sysPostDao := sysdao.NewSysPostDao(ctx)
	data, err := sysPostDao.CheckPostCodeUnique(post.PostCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if data.PostID != 0 && data.PostID != post.PostID {
		return false
	}
	return true
}

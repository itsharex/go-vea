package syssrv

import (
	"context"
	"go-vea/app/common/page"
	"go-vea/app/dao/sysdao"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
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
	err := sysPostDao.Insert(sysPost)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return nil
}

func (*SysPostService) UpdatePostById(ctx context.Context, sysPost *system.SysPost) error {
	sysPostDao := sysdao.NewSysPostDao(ctx)
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

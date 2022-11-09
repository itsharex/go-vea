package syssrv

import (
	"context"
	"go-web-template/app/common/page"
	"go-web-template/app/dao/sysdao"
	"go-web-template/app/model/system"
	"go-web-template/app/model/system/request"
	"go-web-template/global"
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

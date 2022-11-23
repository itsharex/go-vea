package monitorsrv

import (
	"context"
	"go-vea/app/common/page"
	"go-vea/app/dao/monitordao"
	"go-vea/app/model/monitor"
	"go-vea/app/model/monitor/request"
	"go-vea/global"
)

type SysOperLogService struct{}

var SysOperLogSrv = new(SysOperLogService)

func (*SysOperLogService) SelectSysOperLogList(ctx context.Context, operLog *request.SysOperLog) (*page.Pagination, error) {
	sysOperLogDao := monitordao.NewSysOperLogDao(ctx)
	data, err := sysOperLogDao.SelectList(operLog)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysOperLogService) SelectSysOperLog(ctx context.Context, id int64) (*monitor.SysOperLog, error) {
	sysOperLogDao := monitordao.NewSysOperLogDao(ctx)
	data, err := sysOperLogDao.SelectById(id)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysOperLogService) AddSysOperLog(ctx context.Context, operLog *monitor.SysOperLog) error {
	sysOperLogDao := monitordao.NewSysOperLogDao(ctx)
	err := sysOperLogDao.Insert(operLog)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

func (*SysOperLogService) UpdateSysOperLog(ctx context.Context, operLog *monitor.SysOperLog) error {
	sysOperLogDao := monitordao.NewSysOperLogDao(ctx)
	err := sysOperLogDao.UpdateById(operLog)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

func (*SysOperLogService) DeleteSysOperLog(ctx context.Context, ids []int64) error {
	sysOperLogDao := monitordao.NewSysOperLogDao(ctx)
	err := sysOperLogDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

func (*SysOperLogService) CleanOperLog(ctx context.Context) error {
	sysOperLogDao := monitordao.NewSysOperLogDao(ctx)
	err := sysOperLogDao.CleanOperLog()
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

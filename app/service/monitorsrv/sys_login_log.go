package monitorsrv

import (
	"context"
	"go-vea/app/common/page"
	"go-vea/app/dao/monitordao"
	"go-vea/app/model/monitor"
	"go-vea/app/model/monitor/request"
	"go-vea/global"
)

type SysLoginLogService struct{}

var SysLoginLogSrv = new(SysLoginLogService)

func (*SysLoginLogService) SelectSysLoginLogList(ctx context.Context, sysLoginLog *request.SysLoginLog) (*page.Pagination, error) {
	sysLoginLogDao := monitordao.NewSysLoginLogDao(ctx)
	data, err := sysLoginLogDao.SelectList(sysLoginLog)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysLoginLogService) SelectSysLoginLog(ctx context.Context, id int64) (*monitor.SysLoginLog, error) {
	sysLoginLogDao := monitordao.NewSysLoginLogDao(ctx)
	data, err := sysLoginLogDao.SelectById(id)
	if err != nil {
		global.Logger.Error(err)
		return data, err
	}
	return data, err
}

func (*SysLoginLogService) AddSysLoginLog(ctx context.Context, sysLoginLog *monitor.SysLoginLog) error {
	sysLoginLogDao := monitordao.NewSysLoginLogDao(ctx)
	err := sysLoginLogDao.Insert(sysLoginLog)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

func (*SysLoginLogService) UpdateSysLoginLog(ctx context.Context, sysLoginLog *monitor.SysLoginLog) error {
	sysLoginLogDao := monitordao.NewSysLoginLogDao(ctx)
	err := sysLoginLogDao.UpdateById(sysLoginLog)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

func (*SysLoginLogService) DeleteSysLoginLog(ctx context.Context, ids []int64) error {
	sysLoginLogDao := monitordao.NewSysLoginLogDao(ctx)
	err := sysLoginLogDao.DeleteByIds(ids)
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

func (*SysLoginLogService) CleanLoginLog(ctx context.Context) error {
	sysLoginLogDao := monitordao.NewSysLoginLogDao(ctx)
	err := sysLoginLogDao.CleanLoginLog()
	if err != nil {
		global.Logger.Error(err)
		return err
	}
	return err
}

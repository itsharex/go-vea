package monitorctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/framework"
	"go-vea/app/model/monitor"
	"go-vea/app/model/monitor/request"
	"go-vea/app/service/monitorsrv"
	"strconv"
)

type SysLoginLog struct{}

func (*SysLoginLog) GetLoginLogList(ctx *gin.Context) {
	var params request.SysLoginLog
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := monitorsrv.SysLoginLogSrv.SelectSysLoginLogList(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysLoginLog) GetLoginLog(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := monitorsrv.SysLoginLogSrv.SelectSysLoginLog(ctx.Request.Context(), int64(id))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysLoginLog) AddLoginLog(ctx *gin.Context) {
	var params monitor.SysLoginLog
	_ = ctx.ShouldBindJSON(&params)
	err := monitorsrv.SysLoginLogSrv.AddSysLoginLog(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysLoginLog) UpdateLoginLog(ctx *gin.Context) {
	var params monitor.SysLoginLog
	_ = ctx.ShouldBindJSON(&params)
	err := monitorsrv.SysLoginLogSrv.UpdateSysLoginLog(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysLoginLog) DeleteLoginLog(ctx *gin.Context) {
	var params request.SysLoginLog
	_ = ctx.ShouldBindJSON(&params)
	err := monitorsrv.SysLoginLogSrv.DeleteSysLoginLog(ctx.Request.Context(), params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysLoginLog) CleanLoginLog(ctx *gin.Context) {
	err := monitorsrv.SysLoginLogSrv.CleanLoginLog(ctx.Request.Context())
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysLoginLog) Unlock(ctx *gin.Context) {
	framework.SysPasswordSrv.ClearLoginRecordCache(ctx.Param("username"))
	//if err != nil {
	//	result.FailWithMessage(err.Error(), ctx)
	//} else {
	//	result.Ok(ctx)
	//}
	result.Ok(ctx)
}

package monitorctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/model/monitor"
	"go-vea/app/model/monitor/request"
	"go-vea/app/service/monitorsrv"
	"strconv"
)

type SysOperLog struct{}

func (*SysOperLog) GetOperLogList(ctx *gin.Context) {
	var params request.SysOperLog
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := monitorsrv.SysOperLogSrv.SelectSysOperLogList(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysOperLog) GetOperLog(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := monitorsrv.SysOperLogSrv.SelectSysOperLog(ctx.Request.Context(), int64(id))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysOperLog) AddOperLog(ctx *gin.Context) {
	var params monitor.SysOperLog
	_ = ctx.ShouldBindJSON(&params)
	err := monitorsrv.SysOperLogSrv.AddSysOperLog(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysOperLog) UpdateOperLog(ctx *gin.Context) {
	var params monitor.SysOperLog
	_ = ctx.ShouldBindJSON(&params)
	err := monitorsrv.SysOperLogSrv.UpdateSysOperLog(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysOperLog) DeleteOperLog(ctx *gin.Context) {
	var params request.SysOperLog
	_ = ctx.ShouldBindJSON(&params)
	err := monitorsrv.SysOperLogSrv.DeleteSysOperLog(ctx.Request.Context(), params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysOperLog) CleanOperLog(ctx *gin.Context) {
	err := monitorsrv.SysOperLogSrv.CleanOperLog(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysNoticeApi struct{}

func (*SysNoticeApi) GetSysNoticeList(ctx *gin.Context) {
	var params request.SysNotice
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysNoticeSrv.GetSysNoticeList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysNoticeApi) GetSysNotice(ctx *gin.Context) {
	noticeId, _ := strconv.Atoi(ctx.Param("noticeId"))
	data, err := syssrv.SysNoticeSrv.GetSysNoticeById(ctx, int64(noticeId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysNoticeApi) AddSysNotice(ctx *gin.Context) {
	var params system.SysNotice
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysNoticeSrv.AddSysNotice(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysNoticeApi) UpdateSysNotice(ctx *gin.Context) {
	var params system.SysNotice
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysNoticeSrv.UpdateNoticeById(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysNoticeApi) DeleteSysNotice(ctx *gin.Context) {
	var params request.SysNotice
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysNoticeSrv.DeleteSysNoticeByIds(ctx, params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

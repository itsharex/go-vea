package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysDictTypeApi struct{}

func (*SysDictTypeApi) GetDictTypeList(ctx *gin.Context) {
	var params request.SysDictType
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysDictTypeSrv.SelectDictTypeList(ctx.Request.Context(), params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDictTypeApi) GetDictType(ctx *gin.Context) {
	dictCode, _ := strconv.Atoi(ctx.Param("dictId"))
	data, err := syssrv.SysDictTypeSrv.SelectDictTypeById(ctx.Request.Context(), int64(dictCode))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDictTypeApi) AddDictType(ctx *gin.Context) {
	var params system.SysDictType
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysDictTypeSrv.AddDictType(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDictTypeApi) UpdateDictType(ctx *gin.Context) {
	var params system.SysDictType
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysDictTypeSrv.UpdateDictType(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDictTypeApi) DeleteDictType(ctx *gin.Context) {
	var params request.SysDictType
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysDictTypeSrv.DeleteDictTypeByIds(ctx.Request.Context(), params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDictTypeApi) RefreshCache(ctx *gin.Context) {
	err := syssrv.SysDictTypeSrv.ResetDictCache(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDictTypeApi) OptionSelect(ctx *gin.Context) {
	data, err := syssrv.SysDictTypeSrv.SelectDictTypeAll(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

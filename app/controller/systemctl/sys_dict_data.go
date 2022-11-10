package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysDictDataApi struct{}

func (*SysDictDataApi) GetDictDataList(ctx *gin.Context) {
	var params request.SysDictData
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysDictDataSrv.SelectDictDataList(ctx.Request.Context(), params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDictDataApi) GetDictDataListByDictType(ctx *gin.Context) {
	var params request.SysDictData
	params.OpenPage = false
	params.DictType = ctx.Param("dictType")
	data, err := syssrv.SysDictDataSrv.SelectDictDataList(ctx.Request.Context(), params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDictDataApi) GetDictData(ctx *gin.Context) {
	dictCode, _ := strconv.Atoi(ctx.Param("dictCode"))
	data, err := syssrv.SysDictDataSrv.SelectDictDataById(ctx.Request.Context(), int64(dictCode))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDictDataApi) AddDictData(ctx *gin.Context) {
	var params system.SysDictData
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysDictDataSrv.AddDictData(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDictDataApi) UpdateDictData(ctx *gin.Context) {
	var params system.SysDictData
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysDictDataSrv.UpdateDictData(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDictDataApi) DeleteDictData(ctx *gin.Context) {
	var params request.SysDictData
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysDictDataSrv.DeleteDictDataByIds(ctx.Request.Context(), params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

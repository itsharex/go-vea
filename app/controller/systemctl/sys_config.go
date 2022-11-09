package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-web-template/app/common/result"
	"go-web-template/app/model/system"
	"go-web-template/app/model/system/request"
	"go-web-template/app/service/syssrv"
	"strconv"
)

type SysConfigApi struct{}

func (s *SysConfigApi) GetSysConfigList(ctx *gin.Context) {
	var params request.SysConfig
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysConfigSrv.SelectSysConfigList(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (s *SysConfigApi) GetSysConfigById(ctx *gin.Context) {
	configId, _ := strconv.Atoi(ctx.Param("configId"))
	data, err := syssrv.SysConfigSrv.SelectSysConfigById(ctx.Request.Context(), int64(configId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (s *SysConfigApi) GetSysConfigByKey(ctx *gin.Context) {
	data, err := syssrv.SysConfigSrv.SelectSysConfigByKey(ctx.Request.Context(), ctx.Param("configKey"))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (s *SysConfigApi) AddSysConfig(ctx *gin.Context) {
	var config system.SysConfig
	_ = ctx.ShouldBindJSON(&config)
	err := syssrv.SysConfigSrv.AddSysConfig(ctx.Request.Context(), &config)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (s *SysConfigApi) UpdateSysConfig(ctx *gin.Context) {
	var config system.SysConfig
	_ = ctx.ShouldBindJSON(&config)
	err := syssrv.SysConfigSrv.UpdateSysConfig(ctx.Request.Context(), &config)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (s *SysConfigApi) DeleteSysConfig(ctx *gin.Context) {
	var config request.SysConfig
	_ = ctx.ShouldBindJSON(&config)
	err := syssrv.SysConfigSrv.DeleteSysConfigByIds(ctx.Request.Context(), config.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

// RefreshCache 刷新配置缓存
func (s *SysConfigApi) RefreshCache(ctx *gin.Context) {
	syssrv.SysConfigSrv.ResetConfigCache(ctx)
}

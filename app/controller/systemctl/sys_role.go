package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-web-template/app/common/result"
	"go-web-template/app/model/system"
	"go-web-template/app/model/system/request"
	"go-web-template/app/service/syssrv"
	"strconv"
)

type SysRoleApi struct{}

func (*SysRoleApi) GetSysRoleList(ctx *gin.Context) {
	var params request.SysRole
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysRoleSrv.GetSysRoleList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysRoleApi) GetSysRole(ctx *gin.Context) {
	roleId, _ := strconv.Atoi(ctx.Param("roleId"))
	data, err := syssrv.SysRoleSrv.GetSysRoleById(ctx, int64(roleId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysRoleApi) AddSysRole(ctx *gin.Context) {
	var params system.SysRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysRoleSrv.AddSysRole(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysRoleApi) UpdateSysRole(ctx *gin.Context) {
	var params system.SysRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysRoleSrv.UpdateRoleById(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		// todo 更新缓存用户权限
		result.Ok(ctx)
	}
}

func (*SysRoleApi) DeleteSysRole(ctx *gin.Context) {
	var params request.SysRole
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysRoleSrv.DeleteSysRoleByIds(ctx, params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

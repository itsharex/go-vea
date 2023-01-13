package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/core"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysMenuApi struct{}

func (*SysMenuApi) GetMenuList(ctx *gin.Context) {
	var params request.SysMenu
	_ = ctx.ShouldBindJSON(&params)
	loginUser, err := core.TokenSrv.GetLoginUser(ctx)
	data, err := syssrv.SysMenuSrv.SelectMenuList(ctx.Request.Context(), &params, loginUser.UserID)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysMenuApi) GetMenuTreeList(ctx *gin.Context) {
	var params request.SysMenu
	_ = ctx.ShouldBindJSON(&params)
	loginUser, err := core.TokenSrv.GetLoginUser(ctx)
	data, err := syssrv.SysMenuSrv.GetTreeSelect(ctx.Request.Context(), &params, loginUser.UserID)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysMenuApi) GetMenuInfo(ctx *gin.Context) {
	menuId, _ := strconv.Atoi(ctx.Param("menuId"))
	data, err := syssrv.SysMenuSrv.SelectSysMenuById(ctx.Request.Context(), int64(menuId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysMenuApi) AddSysMenu(ctx *gin.Context) {
	var params system.SysMenu
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysMenuSrv.AddSysMenu(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysMenuApi) UpdateSysMenu(ctx *gin.Context) {
	var params system.SysMenu
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysMenuSrv.UpdateSysMenuById(ctx.Request.Context(), &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysMenuApi) DeleteSysMenu(ctx *gin.Context) {
	menuId, _ := strconv.Atoi(ctx.Param("menuId"))
	err := syssrv.SysMenuSrv.DeleteSysMenuByIds(ctx.Request.Context(), int64(menuId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysMenuApi) TreeSelect(ctx *gin.Context) {
	var params request.SysMenu
	_ = ctx.ShouldBindJSON(&params)
	loginUser, err := core.TokenSrv.GetLoginUser(ctx)
	data, err := syssrv.SysMenuSrv.GetTreeSelect(ctx.Request.Context(), &params, loginUser.UserID)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysMenuApi) RoleMenuTreeSelect(ctx *gin.Context) {
	var params request.SysMenu
	_ = ctx.ShouldBindJSON(&params)
	loginUser, err := core.TokenSrv.GetLoginUser(ctx)
	menus, err := syssrv.SysMenuSrv.GetTreeSelect(ctx.Request.Context(), &params, loginUser.UserID)

	checkedKeys, err := syssrv.SysMenuSrv.SelectMenuListByRoleId(ctx, params.RoleId)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithDetailed(response.RoleMenuTreeSelectResp{
			CheckedKeys: checkedKeys,
			Menus:       menus,
		}, "查询成功", ctx)
	}
}

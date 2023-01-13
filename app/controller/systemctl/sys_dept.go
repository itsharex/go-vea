package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/core"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysDeptApi struct{}

func (*SysDeptApi) GetSysDeptList(ctx *gin.Context) {
	var params request.SysDept
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysDeptSrv.GetSysDeptList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDeptApi) GetSysDeptTreeList(ctx *gin.Context) {
	var params request.SysDept
	//params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysDeptSrv.GetDeptTreeList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDeptApi) ExcludeChild(ctx *gin.Context) {
	var params request.SysDept
	_ = ctx.ShouldBindJSON(&params)
	deptList, err := syssrv.SysDeptSrv.GetDeptTreeListExcludeChild(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(deptList, ctx)
	}
}

func (*SysDeptApi) GetSysDept(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))
	if !core.CheckSrv.CheckDeptDataScope(ctx, int64(deptId)) {
		result.FailWithMessage("没有权限访问部门数据", ctx)
		return
	}
	data, err := syssrv.SysDeptSrv.GetSysDeptById(ctx, int64(deptId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysDeptApi) AddSysDept(ctx *gin.Context) {
	var params system.SysDept
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysDeptSrv.AddSysDept(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDeptApi) UpdateSysDept(ctx *gin.Context) {
	var params system.SysDept
	_ = ctx.ShouldBindJSON(&params)
	if !core.CheckSrv.CheckDeptDataScope(ctx, params.DeptID) {
		result.FailWithMessage("没有权限访问部门数据", ctx)
		return
	}
	err := syssrv.SysDeptSrv.UpdateDeptById(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysDeptApi) DeleteSysDept(ctx *gin.Context) {
	deptId, _ := strconv.Atoi(ctx.Param("deptId"))
	err := syssrv.SysDeptSrv.DeleteSysDeptById(ctx, int64(deptId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

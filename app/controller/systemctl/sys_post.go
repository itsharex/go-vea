package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysPostApi struct{}

func (*SysPostApi) GetSysPostList(ctx *gin.Context) {
	var params request.SysPost
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysPostSrv.GetSysPostList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysPostApi) GetSysPost(ctx *gin.Context) {
	postId, _ := strconv.Atoi(ctx.Param("postId"))
	data, err := syssrv.SysPostSrv.GetSysPostById(ctx, int64(postId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysPostApi) AddSysPost(ctx *gin.Context) {
	var params system.SysPost
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysPostSrv.AddSysPost(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysPostApi) UpdateSysPost(ctx *gin.Context) {
	var params system.SysPost
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysPostSrv.UpdatePostById(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysPostApi) DeleteSysPost(ctx *gin.Context) {
	var params request.SysPost
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysPostSrv.DeleteSysPostByIds(ctx, params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

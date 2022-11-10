package systemctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/app/service/syssrv"
	"strconv"
)

type SysUserApi struct{}

func (*SysUserApi) GetSysUserList(ctx *gin.Context) {
	var params request.SysUser
	params.OpenPage = true
	_ = ctx.ShouldBindJSON(&params)
	data, err := syssrv.SysUserSrv.GetSysUserList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*SysUserApi) GetSysUser(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))
	data, err := syssrv.SysUserSrv.GetSysUserById(ctx, int64(userId))
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		// todo
		var roles []*system.SysRole
		var posts []*system.SysPost
		var postIds []int64
		var roleIds []int64
		userInfo := response.UserInfoById{
			User:    data,
			Roles:   roles,
			Posts:   posts,
			PostIds: postIds,
			RoleIds: roleIds,
		}
		result.OkWithData(userInfo, ctx)
	}
}

func (*SysUserApi) AddSysUser(ctx *gin.Context) {
	var params system.SysUser
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysUserSrv.AddSysUser(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) UpdateSysUser(ctx *gin.Context) {
	var params system.SysUser
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysUserSrv.UpdateUserById(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) DeleteSysUser(ctx *gin.Context) {
	var params request.SysUser
	_ = ctx.ShouldBindJSON(&params)
	err := syssrv.SysUserSrv.DeleteSysUserByIds(ctx, params.Ids)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

func (*SysUserApi) ResetPwd(ctx *gin.Context) {

}

func (*SysUserApi) ChangeStatus(ctx *gin.Context) {

}

func (*SysUserApi) GetDeptTree(ctx *gin.Context) {
	params := request.SysDept{}
	data, err := syssrv.SysDeptSrv.GetDeptTreeList(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

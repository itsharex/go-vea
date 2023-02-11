package monitorctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/model/system/request"
	"go-vea/app/service/monitorsrv"
	"go-vea/global"
)

type OnlineUserApi struct{}

func (*OnlineUserApi) GetOnlineUser(ctx *gin.Context) {
	var params request.SysUserOnline
	_ = ctx.ShouldBindJSON(&params)
	data, err := monitorsrv.OnlineUserSrv.GetOnlineUser(ctx, &params)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

func (*OnlineUserApi) ForceLogout(ctx *gin.Context) {
	err := global.Redis.Del(ctx, "login_tokens:"+ctx.Param("tokenId")).Err()
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.Ok(ctx)
	}
}

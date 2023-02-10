package monitorctl

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/service/monitorsrv"
)

type ServerInfoApi struct{}

func (*ServerInfoApi) GetServerInfo(ctx *gin.Context) {
	data, err := monitorsrv.ServerInfoSrv.GetServerInfo(ctx)
	if err != nil {
		result.FailWithMessage(err.Error(), ctx)
	} else {
		result.OkWithData(data, ctx)
	}
}

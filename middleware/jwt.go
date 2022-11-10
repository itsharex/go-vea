package middleware

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/e"
	"go-vea/app/common/result"
	"go-vea/app/framework"
	"go-vea/configs"
	"go-vea/global"
	"go-vea/util"
	"strings"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code = 200
		token := getToken(ctx)
		if token == "" {
			code = e.UNAUTHORIZED
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.UNAUTHORIZED
			} else {
				// 二选一
				// 这里以jwt设置的过期时间为准 默认30分钟 不可刷新
				//if time.Now().Unix() > claims.ExpiresAt {
				//	code = e.FORBIDDEN
				//}

				// 这里以redis key的过期时间为准 可刷新
				loginUser, _ := framework.TokenSrv.GetLoginUser(ctx)
				if loginUser != nil {
					global.Logger.Info("jwt过期时间: ", claims.ExpiresAt)
					framework.TokenSrv.VerifyToken(loginUser)
				} else {
					code = e.FORBIDDEN
				}
			}
		}
		if code != e.SUCCESS {
			result.Forbidden(ctx)
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}

func getToken(ctx *gin.Context) string {
	token := ctx.GetHeader(configs.AppConfig.JWT.Header)
	t := strings.Replace(token, "Bearer ", "", 1)
	return t
}

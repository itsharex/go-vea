package middleware

import (
	"github.com/gin-gonic/gin"
	"go-vea/app/common/result"
	"go-vea/app/framework"
	"go-vea/global"
)

// HasPerm 验证用户是否具备某权限
func HasPerm(perms string) gin.HandlerFunc {
	return func(c *gin.Context) {

		loginUser, err := framework.TokenSrv.GetLoginUser(c)

		if err != nil {
			global.Logger.Error(err)
			c.Abort()
			return
		}
		if hasPermissions(loginUser.Permissions, perms) {
			c.Next()
		} else {
			global.Logger.Error("没有权限")
			c.Abort()
			result.Forbidden(c)
			return
		}
	}
}

// HasRole 判断用户是否拥有某个角色
func HasRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role == "" {
			return
		}

		loginUser, _ := framework.TokenSrv.GetLoginUser(c)

		if loginUser == nil || len(loginUser.SysUserResp.Roles) < 0 {
			global.Logger.Warn("没有权限")
			c.Abort()
			result.Forbidden(c)
			return
		}

		for _, sysRole := range loginUser.SysUserResp.SysRoles {
			roleKey := sysRole.RoleKey
			if roleKey == "admin" || role == roleKey {
				c.Next()
				return
			} else {
				global.Logger.Warn("没有权限")
				c.Abort()
				result.Forbidden(c)
				return
			}
		}
	}
}

func hasPermissions(permissions []string, perm string) bool {
	if perm == "*:*:*" {
		return true
	}
	for _, permission := range permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

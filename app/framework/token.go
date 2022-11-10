package framework

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-vea/app/model/system/response"
	"go-vea/configs"
	"go-vea/global"
	"go-vea/util"
	"strings"
	"time"
)

type TokenService struct{}

var TokenSrv = new(TokenService)

var expireTime = 30 * time.Minute

const MillisMinuteTen = 20 * time.Minute

func (t *TokenService) GetLoginUser(ctx *gin.Context) (loginUser *response.LoginUser, err error) {
	token := getToken(ctx)
	if token != "" {
		claims := &util.Claims{}
		claims, err = util.ParseToken(token)
		if err != nil {
			global.Logger.Error(err)
			return nil, err
		}
		userKey := "login_tokens:" + claims.LoginUserKey
		jsonData, _ := global.Redis.Get(ctx, userKey).Result()
		err = json.Unmarshal([]byte(jsonData), &loginUser)
		if err != nil {
			global.Logger.Error(err)
			return nil, err
		}
		return loginUser, nil
	}
	return nil, nil
}

// Logout 退出登录
func (t *TokenService) Logout(ctx *gin.Context) error {
	token := getToken(ctx)
	if token != "" {
		claims, err := util.ParseToken(token)
		if err != nil && claims == nil {
			global.Logger.Error(err)
			return nil
		}
		err = delLoginUser(claims.LoginUserKey)
		// todo 记录登录信息

		return err
	}
	return errors.New("token为空")
}

// SetLoginUser 设置用户身份信息
func (t *TokenService) SetLoginUser(user *response.LoginUser) {
	if user != nil && user.UserKey != "" {
		refreshToken(user)
	}
}

// CreateToken 创建令牌
func (t *TokenService) CreateToken(user *response.LoginUser) (string, error) {
	user.UserKey = uuid.New().String()
	refreshToken(user)

	token, err := util.GenerateToken(user.UserKey)
	if err != nil {
		global.Logger.Error(err, "token签发失败")
	}
	return token, err
}

// VerifyToken 验证令牌有效期，相差不足20分钟，自动刷新缓存
func (t *TokenService) VerifyToken(user *response.LoginUser) {
	et := user.ExpireTime
	s := time.Now().Sub(et)
	if s <= MillisMinuteTen {
		refreshToken(user)
	}
}

func (t *TokenService) SetUserAgent(user *response.LoginUser) {

}

// refreshToken 刷新token
func refreshToken(user *response.LoginUser) {
	user.LoginTime = time.Now()
	user.ExpireTime = user.LoginTime.Add(expireTime)
	// 将用户信息存入redis
	data, err := json.Marshal(&user)
	userKey := "login_tokens:" + user.UserKey
	err = global.Redis.Set(context.Background(), userKey, data, expireTime).Err()
	if err != nil {
		global.Logger.Error(err)
	}
}

func getToken(ctx *gin.Context) string {
	token := ctx.GetHeader(configs.AppConfig.JWT.Header)
	t := strings.Replace(token, "Bearer ", "", 1)
	return t
}

// delLoginUser 删除用户身份信息
func delLoginUser(userKey string) error {
	if userKey != "" {
		uk := "login_tokens:" + userKey
		err := global.Redis.Del(context.Background(), uk).Err()
		return err
	}
	return errors.New("userKey为空")
}

package framework

import (
	"context"
	"errors"
	"go-vea/app/model/system"
	"go-vea/app/model/system/request"
	"go-vea/global"
	"go-vea/util"
	"time"
)

type SysPasswordService struct{}

var SysPasswordSrv = new(SysPasswordService)

// MaxRetryCount todo 优化
const MaxRetryCount int = 5
const LockTime = 10 * time.Minute

func (s *SysPasswordService) Validate(sysUser *system.SysUser, loginBody *request.LoginBody) error {

	retryCount, err := global.Redis.Get(context.Background(), getCacheKey(sysUser.UserName)).Int()
	if err != nil {
		retryCount = 0
	}
	if retryCount >= MaxRetryCount {
		// todo 记录登录日志
		return errors.New("重试次数太多，请稍后再试")
	}
	if !matches(loginBody, sysUser) {
		retryCount = retryCount + 1
		// todo 记录登录日志

		// 重试5次 锁10分钟
		global.Redis.Set(context.Background(), getCacheKey(sysUser.UserName), retryCount, LockTime)
		return errors.New("密码错误")
	} else {
		// 成功 删除重试记录
		clearLoginRecordCache(sysUser.UserName)
	}
	return nil
}

// 密码匹配
// 表单输入的密码与根据用户名查出的密码比对
func matches(loginBody *request.LoginBody, sysUser *system.SysUser) bool {
	return util.PasswordVerify(loginBody.Password, sysUser.Password)
}

func clearLoginRecordCache(loginName string) {
	k := global.Redis.Exists(context.Background(), getCacheKey(loginName))
	if k != nil {
		global.Redis.Del(context.Background(), getCacheKey(loginName))
	}
}

func getCacheKey(username string) string {
	return "pwd_err_cnt:" + username
}

package captcha

import (
	"context"
	"fmt"
	"go-vea/configs"
	"time"
)

var ctx = context.Background()

const CAPTCHA = "captcha:"

type RedisStore struct{}

// Set 实现设置captcha的方法
func (r RedisStore) Set(id string, value string) error {
	key := CAPTCHA + id
	//time.Minute*2：有效时间2分钟
	err := configs.GetRedisClient(ctx).Set(ctx, key, value, time.Minute*2).Err()

	return err
}

// Get 实现获取captcha的方法
func (r RedisStore) Get(id string, clear bool) string {
	key := CAPTCHA + id
	val, err := configs.GetRedisClient(ctx).Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear {
		//clear为true，验证通过，删除这个验证码
		err := configs.GetRedisClient(ctx).Del(ctx, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// Verify 实现验证captcha的方法
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := RedisStore{}.Get(id, clear)
	return v == answer
}

package response

import (
	"go-vea/app/common/base"
	"time"
)

type SysLoginLog struct {
	base.CommonModel
	ID            int64      `json:"id"`            // 访问ID
	Username      string     `json:"username"`      // 用户账号
	Ipaddr        string     `json:"ipaddr"`        // 登录IP地址
	LoginLocation string     `json:"loginLocation"` // 登录地点
	Browser       string     `json:"browser"`       // 浏览器类型
	Os            string     `json:"os"`            // 操作系统
	Status        string     `json:"status"`        // 登录状态（0成功 1失败）
	Msg           string     `json:"msg"`           // 提示消息
	LoginTime     *time.Time `json:"loginTime"`     // 访问时间
}

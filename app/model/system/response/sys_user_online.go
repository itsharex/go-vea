package response

import "time"

type SysUserOnline struct {
	TokenId       string    `json:"tokenId"`       // 会话编号
	DeptName      string    `json:"deptName"`      // 部门名称
	Username      string    `json:"username"`      // 用户名称
	Ipaddr        string    `json:"ipaddr"`        // 登录IP地址
	LoginLocation string    `json:"loginLocation"` // 登录地址
	Browser       string    `json:"browser"`       // 浏览器类型
	Os            string    `json:"os"`            // 操作系统
	LoginTime     time.Time `json:"loginTime"`     // 登录时间
}

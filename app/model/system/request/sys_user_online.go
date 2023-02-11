package request

import "go-vea/app/common/base"

type SysUserOnline struct {
	base.CommonModel
	TokenId  string `json:"tokenId"`  // 会话编号
	Username string `json:"username"` // 用户名称
	Ipaddr   string `json:"ipaddr"`   // 登录IP地址
}

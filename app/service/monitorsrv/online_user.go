package monitorsrv

import (
	"context"
	"encoding/json"
	"go-vea/app/common/page"
	"go-vea/app/model/system/request"
	"go-vea/app/model/system/response"
	"go-vea/global"
)

type OnlineUserService struct{}

var OnlineUserSrv = new(OnlineUserService)

func (*OnlineUserService) GetOnlineUser(ctx context.Context, userOnline *request.SysUserOnline) (*page.Pagination, error) {
	var sysOnlineUserList []*response.SysUserOnline
	keys, err := global.Redis.Keys(ctx, "login_tokens:*").Result()
	for _, key := range keys {
		loginUser := response.LoginUser{}
		jsonData, _ := global.Redis.Get(ctx, key).Result()
		err = json.Unmarshal([]byte(jsonData), &loginUser)
		if userOnline.Username != "" && userOnline.Ipaddr != "" {
			if userOnline.Ipaddr == loginUser.IpAddr && userOnline.Username == loginUser.SysUserResp.SysUser.Username {
				sysOnlineUserList = append(sysOnlineUserList, selectOnlineByInfo(userOnline.Ipaddr, userOnline.Username, loginUser))
			}
		} else if userOnline.Ipaddr != "" {
			if userOnline.Ipaddr == loginUser.IpAddr {
				sysOnlineUserList = append(sysOnlineUserList, selectOnlineByIpaddr(userOnline.Ipaddr, loginUser))
			}
		} else if userOnline.Username != "" && loginUser.SysUserResp.SysUser != nil {
			if userOnline.Username == loginUser.SysUserResp.SysUser.Username {
				sysOnlineUserList = append(sysOnlineUserList, selectOnlineByUserName(userOnline.Username, loginUser))
			}
		} else {
			sysOnlineUserList = append(sysOnlineUserList, loginUserToUserOnline(loginUser))
		}
	}
	data := &page.Pagination{
		Rows: sysOnlineUserList,
	}
	return data, err
}

func selectOnlineByIpaddr(ipaddr string, user response.LoginUser) *response.SysUserOnline {
	if ipaddr == user.IpAddr {
		return loginUserToUserOnline(user)
	}
	return nil
}

func selectOnlineByUserName(username string, user response.LoginUser) *response.SysUserOnline {
	if username == user.SysUserResp.SysUser.Username {
		return loginUserToUserOnline(user)
	}
	return nil
}

func selectOnlineByInfo(ipaddr string, username string, user response.LoginUser) *response.SysUserOnline {
	if ipaddr == user.IpAddr && username == user.SysUserResp.SysUser.Username {
		return loginUserToUserOnline(user)
	}
	return nil
}

func loginUserToUserOnline(user response.LoginUser) *response.SysUserOnline {
	sysOnlineUser := &response.SysUserOnline{
		TokenId:       user.UserKey,
		Username:      user.SysUserResp.SysUser.Username,
		Ipaddr:        user.IpAddr,
		LoginLocation: user.LoginLocation,
		Browser:       user.Browser,
		Os:            user.Os,
		LoginTime:     user.LoginTime,
	}
	if user.SysUserResp.SysDept != nil {
		sysOnlineUser.DeptName = user.SysUserResp.SysDept.DeptName
	}
	return sysOnlineUser
}

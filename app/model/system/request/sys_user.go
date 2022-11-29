package request

import (
	"go-vea/app/common/base"
	"go-vea/app/model/system"
	"time"
)

type SysUser struct {
	base.CommonModel
	UserID      int64     `json:"userId" form:"userId"`     // 用户ID
	DeptID      int64     `json:"deptId" form:"deptId"`     // 部门ID
	RoleID      string    `json:"roleId"`                   // 角色id
	Username    string    `json:"userName" form:"username"` // 用户账号
	Nickname    string    `json:"nickname"`                 // 用户昵称
	UserType    string    `json:"userType"`                 // 用户类型（00系统用户）
	Email       string    `json:"email"`                    // 用户邮箱
	PhoneNumber string    `json:"phoneNumber"`              // 手机号码
	Gender      string    `json:"gender"`                   // 用户性别（0男 1女 2未知）
	Avatar      string    `json:"avatar"`                   // 头像地址
	Password    string    `json:"password"`                 // 密码
	Status      string    `json:"status"`                   // 帐号状态（0正常 1停用）
	DelFlag     string    `json:"delFlag"`                  // 删除标志（0代表存在 2代表删除）
	LoginIP     string    `json:"loginIp"`                  // 最后登录IP
	LoginDate   time.Time `json:"loginDate"`                // 最后登录时间
	CreateBy    string    `json:"createBy"`                 // 创建者
	CreateTime  time.Time `json:"createTime"`               // 创建时间
	UpdateBy    string    `json:"updateBy"`                 // 更新者
	UpdateTime  time.Time ` json:"updateTime"`              // 更新时间
	Remark      string    `json:"remark"`                   // 备注
	DataScope   string    `json:"dataScope"`                // 数据范围
}

type AddSysUser struct {
	SysUser *system.SysUser `json:"sysUser"`
	RoleIds []int64         `json:"roleIds"`
	PostIds []int64         `json:"postIds"`
}

type ResetPwd struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
	UUID     string `json:"uuid"`
}

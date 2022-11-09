package request

import (
	"time"
)

type SysConfig struct {
	OpenPage    bool      `json:"openPage"`                 // 开启分页
	PageNum     int       `json:"pageNum" form:"pageNum"`   // 页码
	PageSize    int       `json:"pageSize" form:"pageSize"` // 每页大小
	Ids         []int64   `json:"ids"`                      // configIds
	ConfigID    int64     `json:"configId"`                 // 参数主键
	ConfigName  string    `json:"configName"`               // 参数名称
	ConfigKey   string    `json:"configKey"`                // 参数键名
	ConfigValue string    `json:"configValue"`              // 参数键值
	ConfigType  string    `json:"configType"`               // 系统内置（Y是 N否）
	CreateBy    string    `json:"createBy"`                 // 创建者
	CreateTime  time.Time `json:"createTime"`               // 创建时间
	UpdateBy    string    `json:"updateBy"`                 // 更新者
	UpdateTime  time.Time `json:"updateTime"`               // 更新时间
	Remark      string    `json:"remark"`                   // 备注
}

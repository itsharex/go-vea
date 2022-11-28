package response

import (
	"go-vea/app/common/base"
	"time"
)

type SysDictType struct {
	base.CommonModel
	DictID     int64     `json:"dictId"`     // 字典主键
	DictName   string    `json:"dictName"`   // 字典名称
	DictType   string    `json:"dictType"`   // 字典类型
	Status     string    `json:"status"`     // 状态（0正常 1停用）
	CreateBy   string    `json:"createBy"`   // 创建者
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateBy   string    `json:"updateBy"`   // 更新者
	UpdateTime time.Time `json:"updateTime"` // 更新时间
	Remark     string    `json:"remark"`     // 备注
}

type SysDictData struct {
	base.CommonModel
	DictCode   int64     `json:"dictCode"`   // 字典编码
	DictSort   int64     `json:"dictSort"`   // 字典排序
	DictLabel  string    `json:"dictLabel"`  // 字典标签
	DictValue  string    `json:"dictValue"`  // 字典键值
	DictType   string    `json:"dictType"`   // 字典类型
	CSSClass   string    `json:"cssClass"`   // 样式属性（其他样式扩展）
	ListClass  string    `json:"listClass"`  // 表格回显样式
	IsDefault  string    `json:"isDefault"`  // 是否默认（Y是 N否）
	Status     string    `json:"status"`     // 状态（0正常 1停用）
	CreateBy   string    `json:"createBy"`   // 创建者
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateBy   string    `json:"updateBy"`   // 更新者
	UpdateTime time.Time `json:"updateTime"` // 更新时间
	Remark     string    `json:"remark"`     // 备注
}

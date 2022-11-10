package response

import (
	"go-vea/app/common/base"
	"time"
)

type SysPost struct {
	base.CommonModel
	PostID     int64     `json:"postId"`     // 岗位ID
	PostCode   string    `json:"postCode"`   // 岗位编码
	PostName   string    `json:"postName"`   // 岗位名称
	PostSort   int64     `json:"post_Sort"`  // 显示顺序
	Status     string    `json:"status"`     // 状态（0正常 1停用）
	CreateBy   string    `json:"createBy"`   // 创建者
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateBy   string    `json:"updateBy"`   // 更新者
	UpdateTime time.Time `json:"updateTime"` // 更新时间
	Remark     string    `json:"remark"`     // 备注
}

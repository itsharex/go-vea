package response

import (
	"go-web-template/app/common/base"
	"time"
)

type SysNotice struct {
	base.CommonModel
	NoticeID      int64     `json:"noticeId"`      // 公告ID
	NoticeTitle   string    `json:"noticeTitle"`   // 公告标题
	NoticeType    string    `json:"noticeType"`    // 公告类型（1通知 2公告）
	NoticeContent []byte    `json:"noticeContent"` // 公告内容
	Status        string    `json:"status"`        // 公告状态（0正常 1关闭）
	CreateBy      string    `json:"createBy"`      // 创建者
	CreateTime    time.Time `json:"createTime"`    // 创建时间
	UpdateBy      string    `json:"updateBy"`      // 更新者
	UpdateTime    time.Time `json:"updateTime"`    // 更新时间
	Remark        string    `json:"remark"`        // 备注
}

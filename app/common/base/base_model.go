package base

type CommonModel struct {
	OpenPage bool    `json:"openPage"`                 // 开启分页
	PageNum  int     `json:"pageNum" form:"pageNum"`   // 页码
	PageSize int     `json:"pageSize" form:"pageSize"` // 每页大小
	Ids      []int64 `json:"ids"`                      // Ids
}

package page

type Pagination struct {
	PageNum    int         `json:"pageNum,omitempty" form:"pageNum"`
	PageSize   int         `json:"pageSize,omitempty" form:"pageSize"`
	Sort       string      `json:"sort,omitempty" form:"sort"`
	TotalPages int         `json:"total_pages,omitempty"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Total      int64       `json:"total"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

func (p *Pagination) GetPage() int {
	if p.PageNum == 0 {
		p.PageNum = 1
	}
	return p.PageNum
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "createTime desc"
	}
	return p.Sort
}

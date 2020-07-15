package paging

type Paging struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (p Paging) Parse() (offset, limit int) {
	offset, limit = Parse(p.Page, p.PageSize)
	return
}

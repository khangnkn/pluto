package paging

const (
	defaultPage     = 1
	defaultPageSize = 50
)

func Parse(page, pageSize int) (offset, limit int) {
	if page < defaultPage {
		page = defaultPage
	}
	if pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}
	offset = (page - defaultPage) * pageSize
	limit = pageSize
	return
}

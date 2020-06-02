package projectapi

type GetProjectParam struct {
	WorkspaceID uint64 `form:"workspace_id" binding:"required"`
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
}

type CreateProjectParams struct {
	Title string `form:"title" json:"title"`
	Desc  string `form:"desc" json:"desc"`
}

type ProjectResponse struct {
	ID           uint64 `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DatasetCount int    `json:"dataset_count"`
	MemberCount  int    `json:"member_count"`
}

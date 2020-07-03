package projectapi

const (
	SrcAllProject = iota + 1
	SrcMyProject
	SrcOtherProject
)

type GetProjectParam struct {
	WorkspaceID uint64 `form:"workspace_id" binding:"required"`
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
	Source      int    `form:"src" binding:"required"`
	UserID      uint64 `form:"userid"`
}

type CreateProjectParams struct {
	WorkspaceID uint64 `form:"workspace_id" json:"workspace_id"`
	Title       string `form:"title" json:"title"`
	Desc        string `form:"desc" json:"desc"`
}

type CreatePermParams struct {
	ProjectID uint64 `form:"-"`
	UserID    uint64 `form:"user_id" binding:"required"`
	Role      int32  `form:"role" binding:"required"`
}

type ProjectResponse struct {
	ID           uint64 `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DatasetCount int    `json:"dataset_count"`
	MemberCount  int    `json:"member_count"`
}

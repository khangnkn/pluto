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

type UpdateProjectRequest struct {
	Title       string `form:"title" json:"title,omitempty"`
	Description string `form:"description" json:"description,omitempty"`
}

type CreatePermParams struct {
	ProjectID uint64 `form:"-"`
	UserID    uint64 `form:"user_id" binding:"required"`
	Role      int32  `form:"role" binding:"required"`
}

type ProjectResponse struct {
	ID             uint64 `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Thumbnail      string `json:"thumbnail"`
	Color          string `json:"color"`
	DatasetCount   int    `json:"dataset_count"`
	MemberCount    int    `json:"member_count"`
	ProjectManager uint64 `json:"project_manager"`
}

type GetProjectResponse struct {
	Total    int               `json:"total"`
	Projects []ProjectResponse `json:"projects"`
}

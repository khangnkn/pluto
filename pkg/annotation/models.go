package annotation

type PushTaskMessage struct {
	Workspace WorkspaceObject `json:"workspace"`
	Project   ProjectObject   `json:"project"`
	Dataset   DatasetObject   `json:"dataset"`
	Tasks     []TaskObject    `json:"tasks"`
	Labels    []LabelObject   `json:"labels"`
}

type WorkspaceObject struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Admin uint64 `json:"admin"`
}

type ProjectObject struct {
	ID             uint64 `json:"id"`
	Title          string `json:"title"`
	ProjectManager uint64 `json:"project_manager"`
}

type DatasetObject struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	ProjectID uint64 `json:"project_id"`
}

type TaskObject struct {
	ID        uint64 `json:"id"`
	Labeler   uint64 `json:"labeler"`
	Reviewer  uint64 `json:"reviewer"`
	CreatedAt int64  `json:"created_at"`
}

type LabelObject struct {
	ID    uint64     `json:"id"`
	Name  string     `json:"name"`
	Color string     `json:"color"`
	Tool  ToolObject `json:"tool"`
}

type ToolObject struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type LabelStatsResponse struct {
	Total int `json:"total"`
}

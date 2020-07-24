package statsapi

type TaskStatsResponse struct {
	Processed int `json:"processed"`
	Total     int `json:"total"`
}

package taskapi

type CreateTaskRequest struct {
	DatasetID uint64 `json:"dataset_id" form:"dataset_id" binding:"required"`
	Assigner  uint64 `json:"assigner" form:"assigner" binding:"required"`
	Labeler   uint64 `json:"labeler" form:"labeler" binding:"required"`
	Reviewer  uint64 `json:"reviewer" form:"reviewer" binding:"required"`
	Quantity  int    `json:"quantity" form:"quantity" binding:"required"`
}

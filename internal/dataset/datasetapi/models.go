package datasetapi

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type GetDatasetRequest struct {
	ProjectID uint64 `json:"project_id" form:"project_id"`
}
type CreateDatasetRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	ProjectID   uint64 `form:"project_id" json:"project_id"`
}

type CloneDatasetRequest struct {
	ProjectID uint64 `form:"project_id"`
}

type DatasetResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectID   uint64 `json:"project_id"`
	ImageCount  int    `json:"image_count"`
	UpdatedAt   int64  `json:"updated_at"`
}

func (r *repository) ToDatasetResponse(d dataset.Dataset) DatasetResponse {
	var total int
	imgs, err := r.imgRepo.GetAllImageByDataset(d.ProjectID)
	if err == nil {
		total = len(imgs)
	}
	return DatasetResponse{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		ProjectID:   d.ProjectID,
		ImageCount:  total,
		UpdatedAt:   clock.UnixMillisecondFromTime(d.UpdatedAt),
	}
}

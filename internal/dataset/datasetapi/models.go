package datasetapi

import (
	"github.com/nkhang/pluto/internal/dataset"
)

type CreateDatasetRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	ProjectID   uint64 `form:"project_id"`
}

type CloneDatasetRequest struct {
	ProjectID uint64 `form:"project_id"`
}

type DatasetResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectID   uint64 `json:"project_id"`
}

func ToDatasetResponse(d dataset.Dataset) DatasetResponse {
	return DatasetResponse{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		ProjectID:   d.ProjectID,
	}
}

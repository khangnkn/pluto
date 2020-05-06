package datasetapi

import (
	"github.com/nkhang/pluto/internal/dataset"
)

type DatasetResponse struct {
	ID          uint64
	Title       string
	Description string
	ProjectID   uint64
}

func ToDatasetResponse(d dataset.Dataset) DatasetResponse {
	return DatasetResponse{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
	}
}

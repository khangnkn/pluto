package datasetapi

import (
	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type CreateDatasetRequest struct {
	Title       string `form:"title" json:"title" binding:"required"`
	Description string `form:"description" json:"description"`
}

type CloneDatasetRequest struct {
	Token string `form:"token" json:"token"`
}

type ParseLinkRequest struct {
	Link string `form:"link" json:"link" binding:"required"`
}

type DatasetResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	ProjectID   uint64 `json:"project_id"`
	ImageCount  int    `json:"image_count"`
	UpdatedAt   int64  `json:"updated_at"`
}

type GetLinkResponse struct {
	DatasetResponse
	Token string `json:"token"`
}

func (r *repository) ToDatasetResponse(d dataset.Dataset) DatasetResponse {
	var total int
	imgs, err := r.imgRepo.GetAllImageByDataset(d.ID)
	if err == nil {
		total = len(imgs)
	}
	return DatasetResponse{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		Thumbnail:   d.Thumbnail,
		ProjectID:   d.ProjectID,
		ImageCount:  total,
		UpdatedAt:   clock.UnixMillisecondFromTime(d.UpdatedAt),
	}
}

func (d DatasetResponse) WithToken(token string) GetLinkResponse {
	return GetLinkResponse{
		DatasetResponse: d,
		Token:           token,
	}
}

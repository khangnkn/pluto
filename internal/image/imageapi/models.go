package imageapi

import (
	"mime/multipart"

	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type ImageRequestQuery struct {
	DatasetID uint64 `form:"dataset_id"`
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
}

type UploadRequest struct {
	FileHeader *multipart.FileHeader `form:"file"`
	DatasetID  uint64                `form:"dataset_id"`
}

type ImageResponse struct {
	ID        uint64 `json:"id"`
	CreatedAt int64  `json:"created_at"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int64  `json:"size"`
}

type Config struct {
	Scheme     string
	Endpoint   string
	BucketName string
}

func ToImageResponse(i image.Image) ImageResponse {
	return ImageResponse{
		ID:        i.ID,
		CreatedAt: clock.UnixMillisecondFromTime(i.CreatedAt),
		Title:     i.Title,
		URL:       i.URL,
		Width:     i.Width,
		Height:    i.Height,
		Size:      i.Size,
	}
}

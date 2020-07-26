package imageapi

import (
	"mime/multipart"

	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type ImageRequestQuery struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type UploadRequest struct {
	FileHeader []*multipart.FileHeader `form:"file"`
}

type GetImageRequest struct {
	ID uint64 `json:"id"`
}

type ImageResponse struct {
	ID        uint64 `json:"id"`
	DatasetID uint64 `json:"dataset_id"`
	CreatedAt int64  `json:"created_at"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int64  `json:"size"`
}

type Config struct {
	Scheme          string
	Endpoint        string
	BucketName      string
	ThumbnailBucket string
	BasePath        string
}

func ToImageResponse(i image.Image) ImageResponse {
	return ImageResponse{
		ID:        i.ID,
		DatasetID: i.DatasetID,
		CreatedAt: clock.UnixMillisecondFromTime(i.CreatedAt),
		Title:     i.Title,
		URL:       i.URL,
		Thumbnail: i.Thumbnail,
		Width:     i.Width,
		Height:    i.Height,
		Size:      i.Size,
	}
}

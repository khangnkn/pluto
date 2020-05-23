package imageapi

import (
	"mime/multipart"

	"github.com/nkhang/pluto/internal/image"
)

type ImageRequestQuery struct {
	DatasetID uint64 `form:"dataset_id"`
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
}

type UpdloadRequest struct {
	FileHeader *multipart.FileHeader `form:"file"`
	Name       string                `form:"name"`
}

type ImageResponse struct {
	ID     uint64
	URL    string
	Width  int
	Height int
}

func ToImageResponse(i image.Image) ImageResponse {
	return ImageResponse{
		ID:     i.ID,
		URL:    i.URL,
		Width:  i.Width,
		Height: i.Height,
	}
}

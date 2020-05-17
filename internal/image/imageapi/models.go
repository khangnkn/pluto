package imageapi

import "github.com/nkhang/pluto/internal/image"

type ImageRequestQuery struct {
	DatasetID uint64 `form:"dataset_id"`
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
}

type ImageResponse struct {
	ID     uint64
	URL    string
	Width  uint32
	Height uint32
}

func ToImageResponse(i image.Image) ImageResponse {
	return ImageResponse{
		ID:     i.ID,
		URL:    i.URL,
		Width:  i.Width,
		Height: i.Height,
	}

}

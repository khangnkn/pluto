package imagefx

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/internal/image"
	"github.com/nkhang/pluto/internal/image/imageapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/gin"
)

func provideImageRepository(db *gorm.DB, cache cache.Cache) image.Repository {
	dbRepo := image.NewDBRepository(db)
	return image.NewRepository(dbRepo, cache)
}

func provideService(r image.Repository) gin.IEngine {
	repository := imageapi.NewRepository(r)
	return imageapi.NewService(repository)
}

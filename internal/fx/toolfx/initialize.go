package toolfx

import (
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/tool"
	"github.com/nkhang/pluto/internal/tool/toolapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/gin"
)

func provideToolRepository(db *gorm.DB, cacheRepo cache.Cache) tool.Repository {
	diskRepo := tool.NewDiskRepository(db)
	repository := tool.NewRepository(diskRepo, cacheRepo)
	return repository
}

func provideToolAPI(r tool.Repository) toolapi.Repository {
	return toolapi.NewRepository(r)
}

func provideToolService(r toolapi.Repository) gin.IEngine {
	return toolapi.NewService(r)
}

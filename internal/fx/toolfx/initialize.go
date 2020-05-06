package toolfx

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/internal/tool"
	"github.com/nkhang/pluto/internal/tool/toolapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/gin"
)

func provideToolRepository(db *gorm.DB, rc redis.UniversalClient) tool.Repository {
	diskRepo := tool.NewDiskRepository(db)
	cacheClient := cache.New(rc)
	repository := tool.NewRepository(diskRepo, cacheClient)
	return repository
}

func provideToolAPI(r tool.Repository) toolapi.Repository {
	return toolapi.NewRepository(r)
}

func provideToolService(r toolapi.Repository) gin.IEngine {
	return toolapi.NewService(r)
}

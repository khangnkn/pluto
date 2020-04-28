package toolapifx

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/internal/tool"
	"github.com/nkhang/pluto/internal/toolapi"
	"github.com/nkhang/pluto/pkg/cache"
)

func provideToolAPI(r toolapi.Repository) toolapi.ToolRepository {
	return toolapi.NewRepository(r)
}

func provideToolRepository(db *gorm.DB, rc redis.UniversalClient) toolapi.Repository {
	diskRepo := tool.NewDiskRepository(db)
	cacheClient := cache.New(rc)
	repository := tool.NewRepository(diskRepo, cacheClient)
	return repository
}

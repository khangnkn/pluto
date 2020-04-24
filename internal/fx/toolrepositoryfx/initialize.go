package toolrepositoryfx

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/nkhang/pluto/internal/tool"
	"github.com/nkhang/pluto/internal/toolapi"
	"github.com/nkhang/pluto/pkg/cache"
)

func initializer(db *gorm.DB, rc redis.UniversalClient) toolapi.Repository {
	db.AutoMigrate(&tool.Tool{})
	diskRepo := tool.NewDiskRepository(db)
	cacheClient := cache.New(rc)
	repository := tool.NewRepository(diskRepo, cacheClient)
	return repository
}

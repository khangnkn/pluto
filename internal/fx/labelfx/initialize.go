package labelfx

import (
	"github.com/jinzhu/gorm"

	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/label/labelapi"
	"github.com/nkhang/pluto/pkg/cache"
	"github.com/nkhang/pluto/pkg/pgin"
)

func provideRepository(db *gorm.DB, c cache.Cache) label.Repository {
	db.LogMode(true)
	dbRepo := label.NewDiskRepository(db)
	return label.NewRepository(dbRepo, c)
}

func provideService(r label.Repository) pgin.StandaloneRouter {
	repository := labelapi.NewRepository(r)
	return labelapi.NewService(repository)
}

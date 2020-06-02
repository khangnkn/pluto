package dbfx

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"

	"github.com/nkhang/pluto/pkg/logger"
)

func initializer() (*gorm.DB, error) {
	var dialect = viper.GetString("database.dialect")
	logger.Infof("dialect %s", dialect)
	var url = viper.GetString("database.url")
	logger.Infof("url %s", url)
	return gorm.Open(dialect, url)
}

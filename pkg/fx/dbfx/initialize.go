package dbfx

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

func initializer() (*gorm.DB, error) {
	v := viper.Sub("database")
	var dialect = v.GetString("dialect")
	var url = v.GetString("url")
	return gorm.Open(dialect, url)
}

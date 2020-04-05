package configfx

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func Initialize(filename string) fx.Option {
	return fx.Invoke(func() {
		viper.SetConfigName(filename)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config/")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	})

}
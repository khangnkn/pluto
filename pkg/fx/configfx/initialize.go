package configfx

import (
	"log"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/nkhang/pluto/pkg/logger"
)

func Initialize(filename string) fx.Option {
	return fx.Invoke(func() {
		viper.AutomaticEnv()
		viper.SetEnvPrefix(filename)
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.AllowEmptyEnv(true)
		log.Printf(viper.GetString("database.url"))

		viper.SetConfigName(filename)
		viper.SetConfigType("yaml")
		// viper.AddConfigPath(".")
		viper.AddConfigPath("./config/")
		err := viper.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Print("no configuration file found")
			} else {
				panic(err)
			}
		}
		log.Print("reading prod")
		prod := viper.GetBool("service.production")
		log.Print(prod)
		logger.Initlialize(prod)
	})

}

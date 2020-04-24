package redisfx

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func initializer() (redis.UniversalClient, error) {
	addr := viper.GetString("redis.address")
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	err := client.Ping().Err()
	return client, err
}

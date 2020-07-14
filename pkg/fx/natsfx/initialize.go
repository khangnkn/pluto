package natsfx

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func provideNATSClient(lc fx.Lifecycle) *nats.EncodedConn {
	url := viper.GetString("nats.url")
	logger.Infof("Opening connection to NATS server at %s", url)
	nc, err := nats.Connect(url)
	if err != nil {
		logger.Panic(err)
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		logger.Panic(err)
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			nc.Close()
			return nil
		},
	})
	return ec

}

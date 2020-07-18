package eurekafx

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/spf13/cast"

	"go.uber.org/fx"

	"github.com/spf13/viper"

	"github.com/nkhang/pluto/pkg/logger"

	"github.com/hudl/fargo"
	uuid "github.com/satori/go.uuid"
)

func initialize(lc fx.Lifecycle) fargo.EurekaConnection {
	conn := fargo.NewConn(viper.GetString("eureka.address"))
	ip, err := getIP()
	if err != nil {
		panic(err)
	}
	logger.Info(ip)
	host := viper.GetString("eureka.hostname")
	port := viper.GetInt("service.port")
	app := viper.GetString("eureka.app")
	instanceId := uuid.NewV4().String() + ":" + app + ":" + cast.ToString(port)
	ins := fargo.Instance{
		InstanceId:        instanceId,
		HostName:          host,
		App:               app,
		IPAddr:            ip,
		Status:            fargo.UP,
		Port:              port,
		PortEnabled:       true,
		SecurePort:        8443,
		SecurePortEnabled: false,
		HomePageUrl:       fmt.Sprintf("http://%s:%d/", host, port),
		StatusPageUrl:     fmt.Sprintf("http://%s:%d/status", host, port),
		HealthCheckUrl:    fmt.Sprintf("http://%s:%d/healthcheck", host, port),
		DataCenterInfo: fargo.DataCenterInfo{
			Name: fargo.MyOwn,
		},
		Metadata: fargo.InstanceMetadata{
			Raw: []byte("\"instanceId\":\"vendor:" + instanceId + "\""),
		},
	}
	err = conn.RegisterInstance(&ins)
	if err != nil {
		panic(err)
	}
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return conn.DeregisterInstance(&ins)
			},
		})
	return conn
}

func getIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				return "", errors.New("invalid type")
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("cannot get ip")
}

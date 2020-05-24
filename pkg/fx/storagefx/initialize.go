package storagefx

import (
	"github.com/coreos/etcd/clientv3/balancer/resolver/endpoint"
	"github.com/nkhang/pluto/pkg/objectstorage"
	"github.com/spf13/viper"
)

func provideObjectStorage() (objectstorage.ObjectStorage, error) {
	v := viper.Sub("minio")
	endpoint := v.GetString("endpoint")
	accessKey := v.GetString("accesskey")
	secretKey := v.GetString("secretkey")
	useSSL := v.GetBool("usessl")
	return objectstorage.NewMinioClient(endpoint, accessKey, secretKey, useSSL)

}
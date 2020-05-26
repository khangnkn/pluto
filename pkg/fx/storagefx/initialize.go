package storagefx

import (
	"github.com/nkhang/pluto/pkg/objectstorage"
	"github.com/spf13/viper"
)

func provideObjectStorage() (objectstorage.ObjectStorage, error) {
	endpoint := viper.GetString("minio.endpoint")
	accessKey := viper.GetString("minio.accesskey")
	secretKey := viper.GetString("minio.secretkey")
	useSSL := viper.GetBool("minio.usessl")
	return objectstorage.NewMinioClient(endpoint, accessKey, secretKey, useSSL)
}
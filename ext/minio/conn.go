package minio

import (
	"fmt"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"github.com/minio/minio-go"
)

var client *minio.Client

func connect() *minio.Client {
	minioConfig := utils.Config.Minio
	glog.Infof(
		"minio connection: user=%s password=******* host=%s port=%d ssl=%t",
		minioConfig.User, minioConfig.Host, minioConfig.Port, minioConfig.SSL,
	)
	client, err := minio.New(
		fmt.Sprintf("%s:%d", minioConfig.Host, minioConfig.Port),
		minioConfig.User, minioConfig.Password, minioConfig.SSL,
	)
	if err != nil {
		glog.Fatalf("Something bad happened during minio connection: %s", err)
	}
	return client
}

func Client() *minio.Client {
	if client == nil {
		client = connect()
	}
	return client
}

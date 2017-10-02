package qingcloud

import (
	"log"

	qcConfig "github.com/yunify/qingcloud-sdk-go/config"
	qcService "github.com/yunify/qingcloud-sdk-go/service"
)

type Config struct {
	AccessKey string
	SecretKey string
}

type QingCloudClient struct {
	config  *Config
	service *qcService.QingCloudService
}

func (c *Config) Client() (*QingCloudClient, error) {
	cfg, err := qcConfig.New(c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, err
	}

	srv, err := qcService.Init(cfg)

	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] QingCloud Client created")

	return &QingCloudClient{
		config:  c,
		service: srv,
	}, nil
}

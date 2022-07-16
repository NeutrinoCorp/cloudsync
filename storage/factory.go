package storage

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/neutrinocorp/cloudsync"
)

var ErrInvalidBlobStorage = errors.New("cloudsync: Invalid blob storage")

type StoreType uint8

const (
	UnknownStore StoreType = iota
	AmazonS3Store
	GoogleDriveStore
	GoogleCloudStore
	AzureBlobStore
)

func NewBlobStorage(cfg cloudsync.Config, t StoreType) (cloudsync.BlobStorage, error) {
	switch t {
	case AmazonS3Store:
		awsCfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.Cloud.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.Cloud.AccessKey,
				cfg.Cloud.SecretKey, "")))
		if err != nil {
			return nil, err
		}
		return NewAmazonS3(s3.NewFromConfig(awsCfg), cfg), nil
	default:
		return nil, ErrInvalidBlobStorage
	}
}

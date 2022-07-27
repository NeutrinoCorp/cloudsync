package storage

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/neutrinocorp/cloudsync"
)

// ErrInvalidBlobStorage the given blob storage type is invalid.
var ErrInvalidBlobStorage = errors.New("cloudsync: Invalid blob storage")

// BlobStoreType a kind of blob storage (Amazon S3, Google Drive, Google Cloud Storage and/or Microsoft Azure Blob
// Storage).
type BlobStoreType uint8

const (
	_ BlobStoreType = iota
	// AmazonS3Store blob storage for Amazon Simple Storage Service (S3).
	AmazonS3Store
	// GoogleDriveStore blob storage for Google Drive.
	GoogleDriveStore
	// GoogleCloudStore blob storage for Google Cloud (GCP) Storage Service.
	GoogleCloudStore
	// AzureBlobStore blob storage for Microsoft Azure Blob Storage Service.
	AzureBlobStore

	AmazonS3Str    = "AMAZON_S3"
	GoogleDriveStr = "GOOGLE_DRIVE"
	GoogleCloudStr = "GCP_STORAGE"
	AzureBlobStr   = "MS_AZURE_BLOB"
)

// BlobStoreMap readable name mapping to BlobStoreType.
var BlobStoreMap = map[string]BlobStoreType{
	AmazonS3Str:    AmazonS3Store,
	GoogleDriveStr: GoogleDriveStore,
	GoogleCloudStr: GoogleCloudStore,
	AzureBlobStr:   AzureBlobStore,
}

// NewBlobStorage allocates a new cloudsync.BlobStorage concrete implementation based on given BlobStoreType.
func NewBlobStorage(cfg cloudsync.Config, storageType string) (cloudsync.BlobStorage, error) {
	switch BlobStoreMap[storageType] {
	case AmazonS3Store:
		var credOpts config.LoadOptionsFunc
		if cfg.Cloud.AccessKey != "" && cfg.Cloud.SecretKey != "" {
			credOpts = config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.Cloud.AccessKey,
				cfg.Cloud.SecretKey, ""))
		}
		awsCfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.Cloud.Region), credOpts)
		if err != nil {
			return nil, err
		}
		return NewAmazonS3(s3.NewFromConfig(awsCfg), cfg), nil
	case GoogleDriveStore:
		// TODO: Add G Drive implementation
		return nil, ErrInvalidBlobStorage
	case GoogleCloudStore:
		// TODO: Add GCP Storage implementation
		return nil, ErrInvalidBlobStorage
	case AzureBlobStore:
		// TODO: Add MS Azure Blob implementation
		return nil, ErrInvalidBlobStorage
	default:
		return nil, ErrInvalidBlobStorage
	}
}

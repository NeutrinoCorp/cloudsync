package storage

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/neutrinocorp/cloudsync"
)

type AmazonS3 struct {
	client   *s3.Client
	bucket   *string
	uploader *manager.Uploader
}

var _ cloudsync.BlobStorage = &AmazonS3{}

func NewAmazonS3(c *s3.Client, cfg cloudsync.Config) *AmazonS3 {
	uploader := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 10
		u.PartSize = 10 * 1024 * 1024 // 10 MiB
	})
	return &AmazonS3{client: c, bucket: &cfg.Cloud.Bucket, uploader: uploader}
}

func (a *AmazonS3) Upload(ctx context.Context, f cloudsync.File) error {
	defer f.Data.Close()

	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:            a.bucket,
		Key:               &f.Key,
		Body:              f.Data,
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *AmazonS3) CheckMod(ctx context.Context, key string, modTime time.Time, size int64) (bool, error) {
	out, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket:            a.bucket,
		Key:               &key,
		IfUnmodifiedSince: &modTime,
	})
	if err != nil && strings.Contains(err.Error(), "api error NotFound: Not Found") {
		return true, nil // if not found, then allow object writing
	} else if err != nil {
		return false, err
	}

	return out.ContentLength < size || out.LastModified.Before(modTime), nil
}

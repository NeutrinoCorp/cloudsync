package storage

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/neutrinocorp/cloudsync"
)

// AmazonS3 Amazon Simple Storage Service (S3) concrete implementation of cloudsync.BlobStorage.
type AmazonS3 struct {
	client   *s3.Client
	bucket   *string
	uploader *manager.Uploader
}

// compile-time interface impl. validation.
var _ cloudsync.BlobStorage = &AmazonS3{}

// NewAmazonS3 allocates a new AmazonS3 instance ready to perform underlying S3 API actions using cloudsync.BlobStorage
// API.
func NewAmazonS3(c *s3.Client, cfg cloudsync.Config) *AmazonS3 {
	uploader := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 10
		u.PartSize = 10 * 1024 * 1024 // 10 MiB
	})
	return &AmazonS3{client: c, bucket: &cfg.Cloud.Bucket, uploader: uploader}
}

func (a *AmazonS3) Upload(ctx context.Context, obj cloudsync.Object) error {
	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:            a.bucket,
		Key:               &obj.Key,
		Body:              obj.Data,
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
	if err == nil {
		return out.ContentLength < size || out.LastModified.Before(modTime), nil
	}
	switch {
	case strings.HasSuffix(err.Error(), "api error NotFound: Not Found"):
		return true, nil // if not found, then allow object writing
	case strings.HasSuffix(err.Error(), "api error Forbidden: Forbidden"):
		return false, cloudsync.ErrFatalStorage
	case strings.HasSuffix(err.Error(), "api error PreconditionFailed: Precondition Failed"):
		return false, nil // object exists, ignore error
	default:
		return false, err
	}
}

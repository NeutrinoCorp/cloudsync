package cloudsync

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func ShutdownUploadWorkers(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	select {
	case <-ctx.Done():
		log.Debug().Msg("cloudsync: Shutting down workers")
		close(objectUploadJobQueue)
		close(objectUploadJobQueueErr)
		objectUploadJobQueue = nil
		objectUploadJobQueueErr = nil
		wg.Done()
	}
}

func ListenAndExecuteUploadJobs(ctx context.Context, storage BlobStorage, wg *sync.WaitGroup) {
	for job := range objectUploadJobQueue {
		go func(startTime time.Time, j File) {
			log.Info().
				Str("object_key", j.Key).
				Msgf("cloudsync: Uploading file")
			err := storage.Upload(ctx, j)
			if err != nil && objectUploadJobQueueErr != nil {
				objectUploadJobQueueErr <- ErrFileUpload{
					Key:    j.Key,
					parent: err,
				}
			}
			log.Info().
				Str("took", time.Since(startTime).String()).
				Str("object_key", j.Key).
				Msgf("cloudsync: Uploaded file")
			wg.Done()
		}(time.Now(), job)
		select {
		case <-ctx.Done():
			break
		default:
		}
	}
}

func ListenUploadErrors(ctx context.Context, cfg Config) {
	for err := range objectUploadJobQueueErr {
		if cfg.Scanner.LogErrors {
			log.
				Err(err).
				Str("parent", err.parent.Error()).
				Msg("cloudsync: File upload failed")
		}
		select {
		case <-ctx.Done():
			break
		default:
		}
	}
}

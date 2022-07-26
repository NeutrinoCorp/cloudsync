package cloudsync

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// ShutdownUploadWorkers closes internal job queues and stores new configuration variables (if required).
func ShutdownUploadWorkers(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	select {
	case <-ctx.Done():
		log.Debug().Msg("cloudsync: Shutting down workers")
		if objectUploadJobQueue != nil {
			close(objectUploadJobQueue)
			objectUploadJobQueue = nil
		}
		if objectUploadJobQueueErr != nil {
			close(objectUploadJobQueueErr)
			objectUploadJobQueueErr = nil
		}
		wg.Done()
	}
}

// ListenAndExecuteUploadJobs waits and executes object upload jobs asynchronously received from internal queues.
//
// Will break listening loop if context was cancelled.
func ListenAndExecuteUploadJobs(ctx context.Context, storage BlobStorage, wg *sync.WaitGroup) {
	for job := range objectUploadJobQueue {
		go func(startTime time.Time, obj Object) {
			defer wg.Done()
			if obj.CleanupFunc != nil {
				defer obj.CleanupFunc()
			}
			log.Info().
				Str("object_key", obj.Key).
				Msg("cloudsync: Uploading file")
			err := storage.Upload(ctx, obj)
			DefaultStats.decreaseUploadJobs()
			if err != nil && objectUploadJobQueueErr != nil {
				objectUploadJobQueueErr <- ErrFileUpload{
					Key:    obj.Key,
					Parent: err,
				}
				return
			}
			log.Info().
				Str("took", time.Since(startTime).String()).
				Str("object_key", obj.Key).
				Uint64("total_upload_jobs", DefaultStats.GetTotalUploadJobs()).
				Uint64("jobs_left", DefaultStats.GetCurrentUploadJobs()).
				Msg("cloudsync: Uploaded file")
		}(time.Now(), job)
	}
}

// ListenUploadErrors waits and performs actions when object upload jobs fail. These errors are sent asynchronously
// through an internal error queue as all internal jobs are scheduled the same way.
//
// Will break listening loop if context was cancelled.
func ListenUploadErrors(cfg Config) {
	for err := range objectUploadJobQueueErr {
		if cfg.Scanner.LogErrors {
			log.
				Err(err).
				Str("parent", err.Parent.Error()).
				Msg("cloudsync: File upload failed")
		}
	}
}

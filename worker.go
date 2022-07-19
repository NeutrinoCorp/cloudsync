package cloudsync

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// ShutdownUploadWorkers closes internal job queues and stores new configuration variables (if required).
func ShutdownUploadWorkers(ctx context.Context, wg *sync.WaitGroup, cfg Config) {
	select {
	case <-ctx.Done():
		log.Debug().Msg("cloudsync: Shutting down workers")
		close(objectUploadJobQueue)
		close(objectUploadJobQueueErr)
		objectUploadJobQueue = nil
		objectUploadJobQueueErr = nil
		go func() {
			if err := SaveConfig(cfg); err != nil {
				log.Warn().Str("error", err.Error()).Msg(" Failed to update configuration file")
			}
			wg.Done()
		}()
		wg.Done()
	}
}

// ListenAndExecuteUploadJobs waits and executes object upload jobs asynchronously received from internal queues.
//
// Will break listening loop if context was cancelled.
func ListenAndExecuteUploadJobs(ctx context.Context, storage BlobStorage, wg *sync.WaitGroup) {
	for job := range objectUploadJobQueue {
		go func(startTime time.Time, j Object) {
			log.Info().
				Str("object_key", j.Key).
				Msgf("cloudsync: Uploading file")
			err := storage.Upload(ctx, j)
			DefaultStats.decreaseUploadJobs()
			if err != nil && objectUploadJobQueueErr != nil {
				objectUploadJobQueueErr <- ErrFileUpload{
					Key:    j.Key,
					Parent: err,
				}
			}
			log.Info().
				Str("took", time.Since(startTime).String()).
				Str("object_key", j.Key).
				Uint64("total_upload_jobs", DefaultStats.GetTotalUploadJobs()).
				Uint64("jobs_left", DefaultStats.GetCurrentUploadJobs()).
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

// ListenUploadErrors waits and performs actions when object upload jobs fail. These errors are sent asynchronously
// through an internal error queue as all internal jobs are scheduled the same way.
//
// Will break listening loop if context was cancelled.
func ListenUploadErrors(ctx context.Context, cfg Config) {
	for err := range objectUploadJobQueueErr {
		if cfg.Scanner.LogErrors {
			log.
				Err(err).
				Str("parent", err.Parent.Error()).
				Msg("cloudsync: File upload failed")
		}
		select {
		case <-ctx.Done():
			break
		default:
		}
	}
}

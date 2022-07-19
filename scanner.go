package cloudsync

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Scanner main component which reads and schedules upload jobs based on the files found on directories specified
// in Config.
type Scanner struct {
	cfg           Config
	baseCtx       context.Context
	baseCtxCancel context.CancelFunc
	startTime     time.Time
	shutdownWg    sync.WaitGroup
}

// NewScanner allocates a new Scanner instance which will use specified Config.
func NewScanner(cfg Config) *Scanner {
	return &Scanner{
		cfg:           cfg,
		baseCtx:       nil,
		baseCtxCancel: nil,
		startTime:     time.Time{},
		shutdownWg:    sync.WaitGroup{},
	}
}

// Start bootstraps and runs internal processes to read files and schedule upload jobs.
func (s *Scanner) Start(store BlobStorage) error {
	s.baseCtx, s.baseCtxCancel = context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	ListenForSysInterruption(&s.shutdownWg, s.baseCtxCancel)
	go ListenAndExecuteUploadJobs(s.baseCtx, store, wg)
	go ListenUploadErrors(s.baseCtx, s.cfg)
	s.shutdownWg.Add(2)
	go ShutdownUploadWorkers(s.baseCtx, &s.shutdownWg, s.cfg)

	s.startTime = time.Now()
	log.Info().Msg(" Starting file upload jobs")
	if err := ScheduleFileUploads(s.baseCtx, s.cfg, wg, store); err != nil {
		return err
	}
	wg.Wait()
	return nil
}

// Shutdown stops all internal process gracefully. Moreover, the shutdown process will stop if the specified
// context was cancelled, avoiding application deadlocks if used with context.WithTimeout() in expense of
// a corrupted shutdown.
func (s *Scanner) Shutdown(ctx context.Context) error {
	s.baseCtxCancel()
	select {
	case <-ctx.Done():
		return nil
	default:
		s.shutdownWg.Wait()
		log.Info().
			Str("took", time.Since(s.startTime).String()).
			Uint64("total_upload_jobs", DefaultStats.GetTotalUploadJobs()).
			Msg(" Completed all file upload jobs")
	}
	return nil
}

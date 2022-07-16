package cloudsync

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Scanner struct {
	cfg           Config
	baseCtx       context.Context
	baseCtxCancel context.CancelFunc
	startTime     time.Time
	shutdownWg    sync.WaitGroup
}

func NewScanner(cfg Config) *Scanner {
	return &Scanner{
		cfg:           cfg,
		baseCtx:       nil,
		baseCtxCancel: nil,
		startTime:     time.Time{},
		shutdownWg:    sync.WaitGroup{},
	}
}

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

func (s *Scanner) Shutdown(ctx context.Context) error {
	s.baseCtxCancel()
	select {
	case <-ctx.Done():
		return nil
	default:
		s.shutdownWg.Wait()
		log.Info().
			Str("took", time.Since(s.startTime).String()).
			Msg(" Completed all file upload jobs")
	}
	return nil
}

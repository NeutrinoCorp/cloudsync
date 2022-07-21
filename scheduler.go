package cloudsync

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

// objectUploadJobQueue queue used by scheduler to trigger object upload jobs executions as background tasks.
var objectUploadJobQueue = make(chan Object)

// objectUploadJobQueueErr queue used by scheduler to perform actions when object upload jobs executions running as
// background tasks fail (i.e. logging errors).
var objectUploadJobQueueErr = make(chan ErrFileUpload)

// ListenForSysInterruption waits and gracefully shuts down internal workers when an external agent sends
// a cancellation signal (e.g. pressing Ctrl+C on shell session running the program).
func ListenForSysInterruption(wg *sync.WaitGroup, cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Debug().
			Uint64("total_upload_jobs", DefaultStats.GetTotalUploadJobs()).
			Uint64("current_upload_jobs", DefaultStats.GetCurrentUploadJobs()).
			Msg("cloudsync: System interruption detected, exiting")
		cancel()
		wg.Wait()
		log.Debug().
			Uint64("corrupted_upload_jobs", DefaultStats.GetCurrentUploadJobs()).
			Msg("cloudsync: Gracefully closed all background tasks after interruption")
		os.Exit(1)
	}()
}

// ScheduleFileUploads traverses a directory tree based on specified configuration (Config.RootDirectory) and
// schedules upload jobs for each file found within all directories (if ScannerConfig.DeepTraversing was set
// as true) or files found in root directory only.
//
// Furthermore, based on specified Config, a traversing process might get skipped if folder is hidden (uses
// '.' prefix character) or object/folder key was specified to be ignored explicitly in Config file.
func ScheduleFileUploads(ctx context.Context, cfg Config, wg *sync.WaitGroup, storage BlobStorage) error {
	return filepath.WalkDir(cfg.RootDirectory, func(path string, d fs.DirEntry, err error) error {
		isHidden := d.Name() != "." && strings.HasPrefix(d.Name(), ".")
		if d.IsDir() && (isHidden || cfg.KeyIsIgnored(d.Name())) {
			return fs.SkipDir
		} else if d.IsDir() || (isHidden && !cfg.Scanner.ReadHidden) || cfg.KeyIsIgnored(d.Name()) {
			return nil // ignore
		}

		rel, err := filepath.Rel(cfg.RootDirectory, path)
		if err != nil && objectUploadJobQueueErr != nil {
			objectUploadJobQueueErr <- ErrFileUpload{
				Key:    d.Name(),
				Parent: err,
			}
			return nil
		}

		wg.Add(1)
		go scheduleFileUpload(scheduleFileUploadArgs{
			ctx:          ctx,
			cfg:          cfg,
			dirEntry:     d,
			path:         path,
			relativePath: rel,
			wg:           wg,
			storage:      storage,
		})
		return nil
	})
}

type scheduleFileUploadArgs struct {
	ctx          context.Context
	cfg          Config
	dirEntry     fs.DirEntry
	path         string
	relativePath string
	wg           *sync.WaitGroup
	storage      BlobStorage
}

// scheduleFileUpload performs the actual scheduling process for an object upload job.
//
// In addition, it adds a prefix specified in ScannerConfig.PartitionID to create a logical partition.
func scheduleFileUpload(args scheduleFileUploadArgs) {
	if args.cfg.Scanner.PartitionID != "" {
		args.relativePath = fmt.Sprintf("%s/%s", args.cfg.Scanner.PartitionID, args.relativePath)
	}
	args.relativePath = strings.ReplaceAll(args.relativePath, "\\", "/")

	info, _ := args.dirEntry.Info()
	wasMod, err := args.storage.CheckMod(args.ctx, args.relativePath, info.ModTime(), info.Size())
	if err != nil && errors.Is(err, ErrFatalStorage) {
		panic(err)
	} else if err != nil || !wasMod {
		args.wg.Done()
		return
	}

	var obj *os.File
	obj, err = os.Open(args.path)
	if err != nil && objectUploadJobQueueErr != nil {
		objectUploadJobQueueErr <- ErrFileUpload{
			Key:    args.dirEntry.Name(),
			Parent: err,
		}
		args.wg.Done()
		return
	}
	if objectUploadJobQueue != nil {
		DefaultStats.increaseUploadJobs()
		objectUploadJobQueue <- Object{
			Key:  args.relativePath,
			Data: obj,
			CleanupFunc: func() error {
				return obj.Close()
			},
		}
	}
}

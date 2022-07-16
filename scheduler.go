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

var objectUploadJobQueue = make(chan File)
var objectUploadJobQueueErr = make(chan ErrFileUpload)

func ListenForSysInterruption(wg *sync.WaitGroup, cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		wg.Wait()
		log.Debug().Msg("cloudsync: Gracefully closed all background tasks after interruption, exiting")
		os.Exit(1)
	}()
}

func ScheduleFileUploads(ctx context.Context, cfg Config, wg *sync.WaitGroup, st BlobStorage) error {
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
		go scheduleFileUpload(ctx, cfg, rel, wg, st, path, d)
		return nil
	})
}

func scheduleFileUpload(ctx context.Context, cfg Config, rel string, wg *sync.WaitGroup, st BlobStorage,
	path string, d fs.DirEntry) {
	if cfg.Scanner.PartitionID != "" {
		rel = fmt.Sprintf("%s/%s", cfg.Scanner.PartitionID, rel)
	}
	rel = strings.ReplaceAll(rel, "\\", "/")

	info, _ := d.Info()
	wasMod, err := st.CheckMod(ctx, rel, info.ModTime(), info.Size())
	if err != nil && errors.Is(err, ErrFatalStorage) {
		panic(err)
	} else if err != nil || !wasMod {
		wg.Done()
		return
	}

	var obj *os.File
	obj, err = os.Open(path)
	if err != nil && objectUploadJobQueueErr != nil {
		objectUploadJobQueueErr <- ErrFileUpload{
			Key:    d.Name(),
			Parent: err,
		}
		wg.Done()
		return
	}
	if objectUploadJobQueue != nil {
		objectUploadJobQueue <- File{
			Key:  rel,
			Data: obj,
		}
	}
}

package cloudsync

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func testShutdownUploadWorkers(t *testing.T) {
	wg := sync.WaitGroup{}
	initRoutines := runtime.NumGoroutine()
	go ListenUploadErrors(Config{Scanner: ScannerConfig{LogErrors: false}})
	require.Equal(t, initRoutines+1, runtime.NumGoroutine())

	objectUploadJobQueueErr <- ErrFileUpload{
		Key:    "",
		Parent: errors.New("testShutdownUploadWorkers: foo error"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	ShutdownUploadWorkers(ctx, &wg)
	wg.Wait()
	runtime.Gosched()
	runtime.GC()
	time.Sleep(time.Millisecond)
	require.Equal(t, initRoutines, runtime.NumGoroutine())
}

func testListenUploadErrors(t *testing.T) {
	initRoutines := runtime.NumGoroutine()
	go ListenUploadErrors(Config{Scanner: ScannerConfig{LogErrors: true}})
	objectUploadJobQueueErr <- ErrFileUpload{
		Key:    "",
		Parent: errors.New("testListenUploadErrors: foo error"),
	}
	require.Equal(t, initRoutines+1, runtime.NumGoroutine())
	close(objectUploadJobQueueErr)
	objectUploadJobQueueErr = nil
	runtime.Gosched()
	runtime.GC()
	time.Sleep(time.Millisecond)
	require.Equal(t, initRoutines, runtime.NumGoroutine())
}

func testListenUpload(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	initRoutines := runtime.NumGoroutine()
	storage := &NoopBlobStorage{UploadErr: nil}
	go ListenUploadErrors(Config{})
	go ListenAndExecuteUploadJobs(context.TODO(), storage, &wg)
	objectUploadJobQueue <- Object{
		Key:  "foo",
		Data: bytes.NewReader([]byte("foo")),
		CleanupFunc: func() error {
			log.Debug().Msg("testListenUpload: cleaning up")
			return nil
		},
	}
	time.Sleep(time.Millisecond)
	storage.UploadErr = errors.New("bar error")
	objectUploadJobQueue <- Object{
		Key:         "bar",
		Data:        bytes.NewReader([]byte("bar")),
		CleanupFunc: nil,
	}
	require.Equal(t, initRoutines+2, runtime.NumGoroutine())
	wg.Wait()
	close(objectUploadJobQueue)
	objectUploadJobQueue = nil
	close(objectUploadJobQueueErr)
	objectUploadJobQueueErr = nil
	runtime.Gosched()
	runtime.GC()
	time.Sleep(time.Millisecond)
	require.Equal(t, initRoutines, runtime.NumGoroutine())
}

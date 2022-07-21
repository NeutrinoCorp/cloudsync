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

func TestUploadWorkers(t *testing.T) {
	// this global upload worker func is required to avoid goroutines deadlocks and nil pointer panics
	testListenUpload(t)
	objectUploadJobQueueErr = make(chan ErrFileUpload, 0)
	testListenUploadErrors(t)
	objectUploadJobQueue = make(chan Object, 0)
	objectUploadJobQueueErr = make(chan ErrFileUpload, 0)
	testShutdownUploadWorkers(t)
}

func testShutdownUploadWorkers(t *testing.T) {
	wg := sync.WaitGroup{}
	require.Equal(t, 2, runtime.NumGoroutine())
	go ListenUploadErrors(Config{Scanner: ScannerConfig{LogErrors: false}})
	require.Equal(t, 3, runtime.NumGoroutine())

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
	require.Equal(t, 2, runtime.NumGoroutine())
}

func testListenUploadErrors(t *testing.T) {
	require.Equal(t, 2, runtime.NumGoroutine())
	go ListenUploadErrors(Config{Scanner: ScannerConfig{LogErrors: true}})
	objectUploadJobQueueErr <- ErrFileUpload{
		Key:    "",
		Parent: errors.New("testListenUploadErrors: foo error"),
	}
	require.Equal(t, 3, runtime.NumGoroutine())
	close(objectUploadJobQueueErr)
	runtime.Gosched()
	runtime.GC()
	time.Sleep(time.Millisecond)
	require.Equal(t, 2, runtime.NumGoroutine())
}

func testListenUpload(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	require.Equal(t, 2, runtime.NumGoroutine())
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
	require.Equal(t, 4, runtime.NumGoroutine())
	wg.Wait()
	close(objectUploadJobQueue)
	close(objectUploadJobQueueErr)
	runtime.Gosched()
	runtime.GC()
	time.Sleep(time.Millisecond)
	require.Equal(t, 2, runtime.NumGoroutine())
}

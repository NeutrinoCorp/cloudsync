package cloudsync

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testScheduleFileUploads(t *testing.T) {
	tests := []struct {
		name            string
		cfg             Config
		storage         BlobStorage
		expReceivedKeys map[string]struct{}
		expErr          bool
	}{
		{
			name: "Arbitrary root dir",
			cfg: Config{
				RootDirectory: "#@Ad",
			},
			storage: NoopBlobStorage{
				CheckModBool: false,
				CheckModErr:  nil,
			},
			expReceivedKeys: map[string]struct{}{},
			expErr:          true,
		},
		{
			name: "Storage critical error",
			cfg: Config{
				RootDirectory: "./testdata",
			},
			storage: NoopBlobStorage{
				CheckModBool: false,
				CheckModErr:  ErrFatalStorage,
			},
			expReceivedKeys: map[string]struct{}{},
		},
		{
			name: "Storage non-critical error",
			cfg: Config{
				RootDirectory: "./testdata",
			},
			storage: NoopBlobStorage{
				CheckModBool: false,
				CheckModErr:  errors.New("foo error"),
			},
			expReceivedKeys: map[string]struct{}{},
		},
		{
			name: "No modified files",
			cfg: Config{
				RootDirectory: "./testdata",
			},
			storage: NoopBlobStorage{
				CheckModBool: false,
			},
			expReceivedKeys: map[string]struct{}{},
		},
		{
			name: "No hidden files",
			cfg: Config{
				RootDirectory: "./testdata",
				Scanner: ScannerConfig{
					PartitionID:    "",
					ReadHidden:     false,
					DeepTraversing: false,
					IgnoredKeys:    nil,
					LogErrors:      false,
				},
			},
			storage: NoopBlobStorage{
				CheckModBool: true,
			},
			expReceivedKeys: map[string]struct{}{
				"config.1.yaml": {},
				"config.2.yaml": {},
				"config.yaml":   {},
			},
		},
		{
			name: "No deep traversing",
			cfg: Config{
				RootDirectory: "./testdata",
				Scanner: ScannerConfig{
					PartitionID:    "",
					ReadHidden:     true,
					DeepTraversing: false,
					IgnoredKeys:    nil,
					LogErrors:      false,
				},
			},
			storage: NoopBlobStorage{
				CheckModBool: true,
			},
			expReceivedKeys: map[string]struct{}{
				".gitkeep":      {},
				"config.1.yaml": {},
				"config.2.yaml": {},
				"config.yaml":   {},
			},
		},
		{
			name: "Ignore folder",
			cfg: Config{
				RootDirectory: "./testdata",
				Scanner: ScannerConfig{
					PartitionID:    "",
					ReadHidden:     true,
					DeepTraversing: true,
					IgnoredKeys:    []string{"foo"},
					LogErrors:      false,
				},
			},
			storage: NoopBlobStorage{
				CheckModBool: true,
			},
			expReceivedKeys: map[string]struct{}{
				".gitkeep":      {},
				"config.1.yaml": {},
				"config.2.yaml": {},
				"config.yaml":   {},
			},
		},
		{
			name: "No sharding",
			cfg: Config{
				RootDirectory: "./testdata",
				Scanner: ScannerConfig{
					PartitionID:    "",
					ReadHidden:     true,
					DeepTraversing: true,
					IgnoredKeys:    nil,
					LogErrors:      false,
				},
			},
			storage: NoopBlobStorage{
				CheckModBool: true,
			},
			expReceivedKeys: map[string]struct{}{
				".gitkeep":      {},
				"config.1.yaml": {},
				"config.2.yaml": {},
				"config.yaml":   {},
				"foo/foo.yaml":  {},
				"foo/bar.yaml":  {},
			},
		},
		{
			name: "Complete config",
			cfg: Config{
				RootDirectory: "./testdata",
				Scanner: ScannerConfig{
					PartitionID:    "123",
					ReadHidden:     true,
					DeepTraversing: true,
					IgnoredKeys:    []string{"bar.yaml"},
					LogErrors:      false,
				},
			},
			storage: NoopBlobStorage{
				CheckModBool: true,
			},
			expReceivedKeys: map[string]struct{}{
				"123/.gitkeep":      {},
				"123/config.1.yaml": {},
				"123/config.2.yaml": {},
				"123/config.yaml":   {},
				"123/foo/foo.yaml":  {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			if objectUploadJobQueue == nil {
				objectUploadJobQueue = make(chan Object, 0)
			}
			if objectUploadJobQueueErr == nil {
				objectUploadJobQueueErr = make(chan ErrFileUpload, 0)
			}
			wg := sync.WaitGroup{}
			mu := sync.Mutex{}
			outKeys := make([]string, 0, len(tt.expReceivedKeys))
			err := ScheduleFileUploads(ctx, tt.cfg, &wg, tt.storage)
			assert.Equal(t, tt.expErr, err != nil)
			go func() {
				for work := range objectUploadJobQueue {
					mu.Lock()
					outKeys = append(outKeys, work.Key)
					mu.Unlock()
					_ = work.CleanupFunc()
					wg.Done()
				}
			}()
			go func() { // required to avoid routine deadlock error
				for range objectUploadJobQueueErr {
				}
			}()
			wg.Wait()
			if objectUploadJobQueue != nil {
				close(objectUploadJobQueue)
				objectUploadJobQueue = nil
			}
			if objectUploadJobQueueErr != nil {
				close(objectUploadJobQueueErr)
				objectUploadJobQueueErr = nil
			}

			require.Len(t, outKeys, len(tt.expReceivedKeys))
			for _, v := range outKeys {
				_, ok := tt.expReceivedKeys[v]
				assert.True(t, ok)
			}
		})
	}
}

type fakeSystemSignal struct {
}

var _ os.Signal = fakeSystemSignal{}

func (f fakeSystemSignal) String() string {
	return "fake signal"
}

func (f fakeSystemSignal) Signal() {}

func testListenForSysInterruption(t *testing.T) {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.TODO())
	c := make(chan os.Signal, 2)
	ListenForSysInterruption(&wg, cancel, c)
	go func() {
		c <- fakeSystemSignal{}
	}()
	select {
	case <-time.After(time.Millisecond * 100):
		t.Fatal("testListenForSysInterruption: timeout reached")
	case <-ctx.Done():
		return
	}
}

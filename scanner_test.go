package cloudsync

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testNewScanner(t *testing.T) {
	tests := []struct {
		name         string
		rootDir      string
		storage      BlobStorage
		wantErrStart bool
		wantErrShut  bool
		shutTimeout  time.Duration
	}{
		{
			name:         "Nil storage",
			rootDir:      "./testdata",
			storage:      nil,
			wantErrStart: true,
		},

		{
			name:         "Arbitrary",
			rootDir:      "$dada#@1#qasasd",
			storage:      NoopBlobStorage{},
			wantErrStart: true,
		},
		{
			name:        "Cancel ctx",
			rootDir:     "./testdata",
			storage:     NoopBlobStorage{},
			shutTimeout: 0,
		},
		{
			name:        "Valid",
			rootDir:     "./testdata",
			storage:     NoopBlobStorage{},
			shutTimeout: time.Millisecond,
		},
	}

	for _, tt := range tests {
		scanner := NewScanner(Config{
			RootDirectory: tt.rootDir,
		})

		require.Equal(t, tt.wantErrStart, scanner.Start(tt.storage) != nil)
		if tt.storage == nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.TODO(), tt.shutTimeout)
		assert.Equal(t, tt.wantErrShut, scanner.Shutdown(ctx) != nil)
		cancel()
	}
}

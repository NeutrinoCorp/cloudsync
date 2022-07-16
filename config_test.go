package cloudsync_test

import (
	"testing"

	"github.com/neutrinocorp/cloudsync"
	"github.com/stretchr/testify/assert"
)

func TestConfig_KeyIsIgnored(t *testing.T) {
	cfg := cloudsync.Config{
		Scanner: cloudsync.ScannerConfig{
			IgnoredKeys: []string{"*.go"},
		},
	}
	out := cfg.KeyIsIgnored("foo.go")
	assert.True(t, out)
}

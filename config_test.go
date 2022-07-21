package cloudsync_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/neutrinocorp/cloudsync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name                      string
		path, file, rootDirectory string
		exp                       cloudsync.Config
		err                       error
	}{
		{
			name:          "Empty",
			path:          "",
			file:          "",
			rootDirectory: "",
			exp:           cloudsync.Config{},
			err:           &os.PathError{},
		},
		{
			name:          "Empty file name",
			path:          "./testdata",
			file:          "",
			rootDirectory: "",
			exp:           cloudsync.Config{},
			err:           errors.New("yaml: input error: read testdata: The handle is invalid"),
		},
		{
			name:          "Empty root dir",
			path:          "./testdata",
			file:          "config.yaml",
			rootDirectory: "",
			exp: cloudsync.Config{
				Cloud: cloudsync.CloudConfig{
					Region:    "us-east-2",
					Bucket:    "ncorp-dev-cloudsync",
					AccessKey: "XXXX",
					SecretKey: "XXXX",
				},
				Scanner: cloudsync.ScannerConfig{
					PartitionID:    "01G82XT3907RASKY2JW8QSZ2RR",
					ReadHidden:     true,
					DeepTraversing: true,
					IgnoredKeys:    []string{"Bar", "config.yaml", "*.go"},
					LogErrors:      true,
				},
			},
			err: nil,
		},
		{
			name:          "Empty partition id",
			path:          "./testdata",
			file:          "config.1.yaml",
			rootDirectory: "/home/ncorp/Documents",
			exp: cloudsync.Config{
				RootDirectory: "/home/ncorp/Documents",
				Cloud: cloudsync.CloudConfig{
					Region:    "us-east-2",
					Bucket:    "ncorp-dev-cloudsync",
					AccessKey: "XXXX",
					SecretKey: "XXXX",
				},
				Scanner: cloudsync.ScannerConfig{
					ReadHidden:     true,
					DeepTraversing: false,
					IgnoredKeys:    []string{"Bar", "config.yaml", "*.go"},
					LogErrors:      true,
				},
			},
			err: nil,
		},
		{
			name:          "Full data",
			path:          "./testdata",
			file:          "config.yaml",
			rootDirectory: "/home/ncorp/Documents",
			exp: cloudsync.Config{
				RootDirectory: "/home/ncorp/Documents",
				Cloud: cloudsync.CloudConfig{
					Region:    "us-east-2",
					Bucket:    "ncorp-dev-cloudsync",
					AccessKey: "XXXX",
					SecretKey: "XXXX",
				},
				Scanner: cloudsync.ScannerConfig{
					PartitionID:    "01G82XT3907RASKY2JW8QSZ2RR",
					ReadHidden:     true,
					DeepTraversing: true,
					IgnoredKeys:    []string{"Bar", "config.yaml", "*.go"},
					LogErrors:      true,
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := cloudsync.NewConfig(tt.path, tt.file, tt.rootDirectory)
			require.IsType(t, tt.err, err)
			assert.Equal(t, tt.exp.Cloud.Region, out.Cloud.Region)
			assert.Equal(t, tt.exp.Cloud.Bucket, out.Cloud.Bucket)
			assert.Equal(t, tt.exp.Cloud.AccessKey, out.Cloud.AccessKey)
			assert.Equal(t, tt.exp.Cloud.SecretKey, out.Cloud.SecretKey)

			if tt.exp.Scanner.PartitionID == "" && err == nil {
				assert.NotEmpty(t, out.Scanner.PartitionID)
			} else {
				assert.Equal(t, tt.exp.Scanner.PartitionID, out.Scanner.PartitionID)
			}
			assert.Equal(t, tt.exp.Scanner.ReadHidden, out.Scanner.ReadHidden)
			assert.Equal(t, tt.exp.Scanner.DeepTraversing, out.Scanner.DeepTraversing)
			assert.Equal(t, tt.exp.Scanner.LogErrors, out.Scanner.LogErrors)
			assert.EqualValues(t, tt.exp.Scanner.IgnoredKeys, out.Scanner.IgnoredKeys)
			assert.EqualValues(t, tt.exp.RootDirectory, out.RootDirectory)
		})
	}
}

func TestConfig_KeyIsIgnored(t *testing.T) {
	tests := []struct {
		name string
		cfg  cloudsync.Config
		key  string
		exp  bool
	}{
		{
			name: "Empty",
			cfg:  cloudsync.Config{},
			key:  "",
			exp:  false,
		},
		{
			name: "Empty key",
			cfg: cloudsync.Config{
				Scanner: cloudsync.ScannerConfig{
					IgnoredKeys: []string{"*.go", "Bar"},
				},
			},
			key: "",
			exp: false,
		},
		{
			name: "Existing key",
			cfg: cloudsync.Config{
				Scanner: cloudsync.ScannerConfig{
					IgnoredKeys: []string{"*.go", "Bar"},
				},
			},
			key: "Bar",
			exp: true,
		},
		{
			name: "Existing wildcard key",
			cfg: cloudsync.Config{
				Scanner: cloudsync.ScannerConfig{
					IgnoredKeys: []string{"*.go", "Bar"},
				},
			},
			key: "foo.go",
			exp: true,
		},
		{
			name: "Non-existing key",
			cfg: cloudsync.Config{
				Scanner: cloudsync.ScannerConfig{
					IgnoredKeys: []string{"*.go", "Bar"},
				},
			},
			key: "baz",
			exp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.cfg.KeyIsIgnored(tt.key)
			assert.Equal(t, tt.exp, out)
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tests := []struct {
		name          string
		cfg           cloudsync.Config
		willCleanup   bool
		willFlushFile bool
		err           error
	}{
		{
			name: "Empty",
			cfg:  cloudsync.Config{},
			err:  &fs.PathError{},
		},
		{
			name: "Non-existent file",
			cfg: cloudsync.Config{
				FilePath: "./testdata/foo.yaml",
			},
			err:         nil,
			willCleanup: true,
		},
		{
			name: "Existent file", // replace/update
			cfg: cloudsync.Config{
				FilePath: "./testdata/config.2.yaml",
			},
			err:           nil,
			willFlushFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cloudsync.SaveConfig(tt.cfg)
			assert.IsType(t, tt.err, err)
			if err != nil {
				return
			}
			if tt.willFlushFile {
				_ = os.Truncate(tt.cfg.FilePath, 0)
			}
			if tt.willCleanup {
				_ = os.Remove(tt.cfg.FilePath) // cleanup
			}
		})
	}
}

func TestSaveConfigIfNotExists(t *testing.T) {
	tests := []struct {
		name            string
		path, file      string
		willCleanupFile bool
		willCleanupPath bool
		exp             bool
	}{
		{
			name: "Empty",
			path: "",
			file: "",
			exp:  false,
		},
		{
			name: "Empty file name",
			path: "./testdata",
			file: "",
			exp:  false,
		},
		{
			name: "Empty path",
			path: "",
			file: "config.n.yaml",
			exp:  false,
		},
		{
			name:            "Existing path",
			path:            "./testdata",
			file:            "conf.test.yaml",
			willCleanupFile: true,
			exp:             true,
		},
		{
			name:            "Non-existing path",
			path:            "./testdata0",
			file:            "conf.test.yaml",
			willCleanupFile: true,
			willCleanupPath: true,
			exp:             true,
		},
		{
			name: "Existing file",
			path: "./testdata",
			file: "config.2.yaml",
			exp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := cloudsync.SaveConfigIfNotExists(tt.path, tt.file)
			assert.Equal(t, tt.exp, out)
			if tt.willCleanupFile {
				_ = os.Remove(filepath.Join(tt.path, tt.file))
			}
			if tt.willCleanupPath {
				_ = os.Remove(tt.path)
			}
		})
	}
}

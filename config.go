package cloudsync

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// CloudConfig remote infrastructure and services configuration.
type CloudConfig struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

// ScannerConfig Scanner configuration.
type ScannerConfig struct {
	// PartitionID a Scanner instance will use this field to create logical partitions in the specified bucket.
	//
	// This could be used in many ways such as:
	//
	// - Create a multi-tenant environment.
	//
	// - Store data from several machines (maybe from within a network) into a single bucket without operational
	// overhead.
	//
	// Note: This field is auto-generated using Unique Lexicographic IDs (ULID) if not found.
	PartitionID string `yaml:"partition_id"`
	// ReadHidden read from files using the '.' character prefix.
	ReadHidden bool `yaml:"read_hidden"`
	// DeepTraversing read every node until leafs are reached from a root directory tree. If set to false,
	// Scanner will read only the root tree files.
	DeepTraversing bool `yaml:"deep_traversing"`
	// IgnoredKeys deny list of custom reserved file or folder keys. Scanner will skip items specified here.
	IgnoredKeys []string `yaml:"ignored_keys"`
	// LogErrors disable or enable logging of errors. Useful for development or overall process visibility purposes.
	LogErrors bool `yaml:"log_errors"`
}

// Config Main application configuration.
type Config struct {
	FilePath      string        `yaml:"-"`
	RootDirectory string        `yaml:"-"`
	Cloud         CloudConfig   `yaml:"cloud"`
	Scanner       ScannerConfig `yaml:"scanner"`

	ignoredKeysHashSet map[string]struct{}
}

// NewConfig allocates a Config instance used by internal components to perform its processes.
func NewConfig(path, file, rootDirectory string) (Config, error) {
	filePath := filepath.Join(path, file)
	f, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.Scanner.PartitionID == "" {
		cfg.Scanner.PartitionID = ulid.Make().String() // set a tenant id by default
	}
	cfg.RootDirectory = rootDirectory
	cfg.FilePath = filePath
	return cfg, nil
}

func newIgnoredKeysSet(keys []string) map[string]struct{} {
	if len(keys) == 0 {
		return nil
	}

	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if strings.HasPrefix(k, "*.") {
			k = strings.TrimPrefix(k, "*.")
		}
		m[k] = struct{}{}
	}
	return m
}

// KeyIsIgnored verifies if a specified key was selected to be ignored.
func (c *Config) KeyIsIgnored(key string) bool {
	if c.ignoredKeysHashSet == nil {
		c.ignoredKeysHashSet = newIgnoredKeysSet(c.Scanner.IgnoredKeys)
	}

	splStr := strings.Split(key, ".")
	if len(splStr) > 1 {
		extension := splStr[len(splStr)-1]
		if _, ok := c.ignoredKeysHashSet[extension]; ok {
			return true
		}
	}
	_, ok := c.ignoredKeysHashSet[key]
	return ok
}

// SaveConfig stores the specified Config into host's physical disk.
func SaveConfig(cfg Config) error {
	log.Debug().Msg("cloudsync: Saving configuration file")
	f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	defer encoder.Close()
	return encoder.Encode(cfg)
}

// SaveConfigIfNotExists creates a path and/or Config file if not found.
//
// If no file was found, it will allocate a ULID as ScannerConfig.PartitionID.
func SaveConfigIfNotExists(path, file string) bool {
	if path == "" || file == "" {
		return false
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(path, os.ModePerm)
	}

	// lookup for config file to verify its existence, stop process if found
	dirTmp, _ := os.ReadDir(path)
	for _, entry := range dirTmp {
		if file == entry.Name() {
			return true // idempotency
		}
	}
	_ = SaveConfig(Config{
		FilePath: filepath.Join(path, file),
		Scanner: ScannerConfig{
			PartitionID: ulid.Make().String(),
		},
	})
	return true
}

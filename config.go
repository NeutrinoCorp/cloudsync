package cloudsync

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/oklog/ulid/v2"

	"gopkg.in/yaml.v3"
)

type CloudConfig struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type ScannerConfig struct {
	PartitionID    string   `yaml:"tenant_id"`
	ReadHidden     bool     `yaml:"read_hidden"`
	DeepTraversing bool     `yaml:"deep_traversing"`
	IgnoredKeys    []string `yaml:"ignored_keys"`
	LogErrors      bool     `yaml:"log_errors"`
}

type Config struct {
	loadedTenantID     bool
	filePath           string
	ignoredKeysHashSet map[string]struct{}

	Cloud   CloudConfig
	Scanner ScannerConfig
}

func NewConfig(path, file string) (Config, error) {
	filePath := filepath.Join(path, file)
	f, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	cfg := Config{
		filePath: filePath,
	}
	if err = decoder.Decode(&cfg); err != nil {
		return cfg, err
	}
	if cfg.Scanner.PartitionID == "" {
		cfg.Scanner.PartitionID = ulid.Make().String() // set a tenant id by default
		cfg.loadedTenantID = true
	}
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

func (c Config) SaveConfig() error {
	if !c.loadedTenantID {
		return nil
	}

	log.Debug().Msg("cloudsync: Updating configuration file")
	f, err := os.OpenFile(c.filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	defer encoder.Close()
	return encoder.Encode(c)
}

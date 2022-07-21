package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/neutrinocorp/cloudsync"
	"github.com/neutrinocorp/cloudsync/storage"
)

func main() {
	var dirName string
	var dirCfg string
	var fileCfg string
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		dirCfg = filepath.Join(homeDir, ".cloudsync")
	}

	flag.StringVar(&dirName, "d", "", "Directory to be scanned")
	flag.StringVar(&dirCfg, "c", dirCfg, "Directory for configuration files")
	flag.StringVar(&fileCfg, "f", "config.yaml", "Configuration file name")
	flag.Parse()

	cloudsync.SaveConfigIfNotExists(dirCfg, fileCfg)
	cfg, err := cloudsync.NewConfig(dirCfg, fileCfg, dirName)
	if err != nil {
		panic(err)
	}

	blobStore, err := storage.NewBlobStorage(cfg, storage.AmazonS3Store)
	if err != nil {
		panic(err)
	}

	scanner := cloudsync.NewScanner(cfg) // blocking I/O
	if err = scanner.Start(blobStore); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err = scanner.Shutdown(ctx); err != nil {
		panic(err)
	}
}

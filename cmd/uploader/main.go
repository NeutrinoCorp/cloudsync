package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/neutrinocorp/cloudsync"
	"github.com/neutrinocorp/cloudsync/storage"
	"github.com/rs/zerolog/log"
)

func main() {
	var dirName string
	var dirCfg string
	var fileCfg string
	var storeType string
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		dirCfg = filepath.Join(homeDir, ".cloudsync")
	}

	flag.StringVar(&dirName, "p", "", "Directory path to be scanned")
	flag.StringVar(&dirCfg, "c", dirCfg, "Directory for configuration files")
	flag.StringVar(&fileCfg, "f", "config.yaml", "Configuration file name")
	flag.StringVar(&storeType, "d", "", "Blob storage driver")
	flag.Parse()

	cloudsync.SaveConfigIfNotExists(dirCfg, fileCfg)
	cfg, err := cloudsync.NewConfig(dirCfg, fileCfg, dirName)
	if err != nil {
		log.Err(err).Msg("Could not load configuration file")
		os.Exit(1)
	}

	blobStore, err := storage.NewBlobStorage(cfg, storeType)
	if err != nil {
		log.Err(err).Msg("Could not load blob storage driver")
		os.Exit(1)
	}

	scanner := cloudsync.NewScanner(cfg) // blocking I/O
	if err = scanner.Start(blobStore); err != nil {
		log.Err(err).Msg("Could not start scanner instance")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err = scanner.Shutdown(ctx); err != nil {
		log.Err(err).Msg("Could not gracefully shutdown scanner instance")
		os.Exit(1)
	}
}

package main

import (
	"context"
	"flag"
	"time"

	"github.com/neutrinocorp/cloudsync"
	"github.com/neutrinocorp/cloudsync/storage"
)

func main() {
	var dirName string
	var dirCfg string
	var fileCfg string
	flag.StringVar(&dirName, "d", "", "Directory to be scanned")
	flag.StringVar(&dirCfg, "c", ".", "Directory for configuration files")
	flag.StringVar(&fileCfg, "cf", "config.yaml", "Configuration file name")
	flag.Parse()

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

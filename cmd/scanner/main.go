package main

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/neutrinocorp/cloudsync"
	"github.com/neutrinocorp/cloudsync/storage"
	"github.com/rs/zerolog/log"
)

func main() {
	var dirName string
	var dirCfg string
	var fileCfg string
	flag.StringVar(&dirName, "d", "", "Directory to be scanned")
	flag.StringVar(&dirCfg, "c", ".", "Directory for configuration files")
	flag.StringVar(&fileCfg, "cf", "config.yaml", "Configuration file name")
	flag.Parse()

	cfg, err := cloudsync.NewConfig(dirCfg, fileCfg)
	if err != nil {
		panic(err)
	}
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Cloud.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.Cloud.AccessKey,
			cfg.Cloud.SecretKey, "")))
	if err != nil {
		panic(err)
	}

	st := storage.NewAmazonS3(s3.NewFromConfig(awsCfg), cfg)
	rootCtx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	shutdownWg := new(sync.WaitGroup)

	cloudsync.ListenForSysInterruption(shutdownWg, cancel)
	go cloudsync.ListenAndExecuteUploadJobs(rootCtx, st, wg)
	go cloudsync.ListenUploadErrors(rootCtx, cfg)
	go cloudsync.ShutdownUploadWorkers(rootCtx, shutdownWg)

	startTime := time.Now()
	log.Info().
		Msg("cloudsync: Starting file upload jobs")
	if err = cloudsync.ScheduleFileUploads(rootCtx, cfg, dirName, wg, st); err != nil {
		panic(err)
	}
	wg.Wait()
	cancel()
	shutdownWg.Wait()
	log.Info().
		Str("took", time.Since(startTime).String()).
		Msg("cloudsync: Completed all file upload jobs")

	if err = cfg.SaveConfig(); err != nil {
		log.Warn().Str("error", err.Error()).Msg("cloudsync: Failed to update configuration file")
	}
}

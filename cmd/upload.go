package cmd

import (
	"context"
	"os"
	"time"

	"github.com/neutrinocorp/cloudsync"
	"github.com/neutrinocorp/cloudsync/storage"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	uploadCmd.Flags().StringP("path", "p", "", "Directory path to be scanned")
	_ = uploadCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload objects from a local directory to a selected blob storage",
	Long: `This command traverses through all objects from specified directory 
and compares its content with the selected blob storage contained items. If the file was modified 
locally, then the command will upload the new object version to blob storage.`,
	TraverseChildren: true,
	Example:          "cloudsync upload -p ./Foo -d AMAZON_S3",
	Run:              upload,
}

func upload(cmd *cobra.Command, _ []string) {
	var dirCfg string
	var fileCfg string
	var dirName string
	var storeType string

	dirCfg, _ = cmd.Flags().GetString("configPath")
	fileCfg, _ = cmd.Flags().GetString("configFile")
	dirName, _ = cmd.Flags().GetString("path")
	storeType, _ = cmd.Flags().GetString("driver")

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

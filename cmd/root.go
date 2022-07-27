package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/neutrinocorp/cloudsync/storage"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudsync [OPTIONS] [COMMANDS]",
	Short: "Neutrino CloudSync CLI is an open-source tool used to upload entire file folders from any host to any cloud",
	Long: `CloudSync CLI is an open-source tool created by Neutrino Corporation used to synchronize local files 
(aka. objects) between a single (up to many) machine(s) and a single (up to many) blob storage(s).

Some of the available blob storages drivers are Amazon Simple Storage Service (S3), Google Cloud Storage, Google Drive and Microsoft Azure Blob Storage.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var dirCfg string
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		dirCfg = filepath.Join(homeDir, ".cloudsync")
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudsync.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly
	rootCmd.PersistentFlags().StringP("configPath", "c", dirCfg, "Directory for configuration files")
	rootCmd.PersistentFlags().StringP("configFile", "f", "config.yaml", "Configuration file name")
	rootCmd.PersistentFlags().StringP("driver", "d", "", "Blob storage driver (available drivers: "+
		strings.Join([]string{storage.AmazonS3Str, storage.GoogleDriveStr, storage.GoogleCloudStr,
			storage.AzureBlobStr}, ", ")+")")

	_ = rootCmd.MarkPersistentFlagRequired("driver")
}

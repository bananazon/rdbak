package cmd

import (
	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:          "backup",
		Aliases:      []string{"b"},
		Short:        "Back your raindrop.io bookmarks up to a YAML file",
		Long:         "Back your raindrop.io bookmarks up to a YAML file",
		PreRun:       backupPreRunCmd,
		Run:          backupRunCmd,
		SilenceUsage: false,
	}
	rd = raindrop.Raindrop{}
)

func init() {
	GetBackupFlags(backupCmd)
	rootCmd.AddCommand(backupCmd)
}

func backupPreRunCmd(cmd *cobra.Command, args []string) {
	rd = *raindrop.New(flagConfigFile, flagPrune, logger)
	err = rd.ParseConfig()
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

func backupRunCmd(cmd *cobra.Command, args []string) {
	err = rd.Backup()
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

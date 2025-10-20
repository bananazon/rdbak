package cmd

import (
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
)

func init() {
	GetBackupFlags(backupCmd)
	rootCmd.AddCommand(backupCmd)
}

func backupPreRunCmd(cmd *cobra.Command, args []string) {

}

func backupRunCmd(cmd *cobra.Command, args []string) {
	err = rd.Backup(flagPrune)
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

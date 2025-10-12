package cmd

import (
	"github.com/gdanko/rdbak/globals"
	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/gdanko/rdbak/util"
	"github.com/spf13/cobra"
)

var (
	backupCmd = &cobra.Command{
		Use:          "backup",
		Short:        "Back your raindrop.io bookmarks up to a JSON file",
		Long:         "Back your raindrop.io bookmarks up to a JSON file",
		PreRunE:      backupPreRunCmd,
		RunE:         backupRunCmd,
		SilenceUsage: false,
	}
	rd = raindrop.Raindrop{}
)

func init() {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, flagNoColor)

	err = globals.SetHomeDirectory()
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	GetBackupFlags(backupCmd)
	rootCmd.AddCommand(backupCmd)
}

func backupPreRunCmd(cmd *cobra.Command, args []string) (err error) {
	rd = *raindrop.New(flagConfigFile, logger)
	err = rd.ParseConfig()
	if err != nil {
		return err
	}

	return nil
}

func backupRunCmd(cmd *cobra.Command, args []string) error {
	err = rd.Backup()
	if err != nil {
		return err
	}

	return nil
}

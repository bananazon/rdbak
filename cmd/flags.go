package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/gdanko/rdbak/pkg/globals"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/spf13/cobra"
)

func GetBackupFlags(cmd *cobra.Command) {
	getCommonFlags(cmd)
	getBackupFlags(cmd)
}

func GetEncryptPasswordFlags(cmd *cobra.Command) {
	getCommonFlags(cmd)
}

func getCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&flagConfigFile, "config", "c", filepath.Join(globals.GetHomeDirectory(), ".config", "rdbak", "config.yaml"), "Specify a config file")
	cmd.Flags().StringVarP(&logLevelStr, "log", "l", defaultLogLevel, fmt.Sprintf("The log level, one of: %s", util.ReturnLogLevels(logLevelMap)))
}

func getBackupFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&flagPrune, "prune", "p", false, "Prune older bookmarks-{timestamp}.yaml files")
}

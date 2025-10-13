package cmd

import (
	"path/filepath"

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
	homeDir, _ := util.GetHome()
	cmd.Flags().StringVarP(&flagConfigFile, "config", "c", filepath.Join(homeDir, ".config", "rdbak", "config.yaml"), "Specify a config file")
}

func getBackupFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&flagPrune, "prune", "p", false, "Prune older bookmarks-{timestamp}.yaml files")
}

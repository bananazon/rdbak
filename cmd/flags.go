package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/gdanko/rdbak/pkg/globals"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/spf13/cobra"
)

func GetBackupFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&flagConfigFile, "config", "c", filepath.Join(globals.GetHomeDirectory(), ".config", "rdbak", "config.yaml"), "Specify a config file")
	cmd.Flags().StringVar(&logLevelStr, "log", defaultLogLevel, fmt.Sprintf("The log level, one of: %s", util.ReturnLogLevels(logLevelMap)))
}

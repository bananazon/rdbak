package cmd

import (
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	defaultLogLevel = "info"
	err             error
	flagConfigFile  string
	flagNoColor     bool
	flagPrune       bool
	homeDir         string
	logger          *logrus.Logger
	rootCmd         = &cobra.Command{
		Use:   "rdbak",
		Short: "rdbak is a command line utility to backup your raindrop.io bookmarks",
		Long:  "rdbak is a command line utility to backup your raindrop.io bookmarks",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			homeDir, _ = util.GetHome()
			logger = util.ConfigureLogger(flagNoColor, homeDir)
		},
	}
	versionFull bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {

}

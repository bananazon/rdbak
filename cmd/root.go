package cmd

import (
	"github.com/gdanko/rdbak/pkg/globals"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	defaultLogLevel = "info"
	err             error
	flagConfigFile  string
	flagNoColor     bool
	logger          *logrus.Logger
	logLevel        logrus.Level
	logLevelStr     string
	logLevelMap     = map[string]logrus.Level{
		"panic": logrus.PanicLevel,
		"fatal": logrus.FatalLevel,
		"error": logrus.ErrorLevel,
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"trace": logrus.TraceLevel,
	}
	rootCmd = &cobra.Command{
		Use:   "rdbak",
		Short: "rdbak is a command line utility to backup your raindrop.io bookmarks",
		Long:  "rdbak is a command line utility to backup your raindrop.io bookmarks",
	}
	versionFull bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, flagNoColor)

	err = globals.SetHomeDirectory()
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}
}

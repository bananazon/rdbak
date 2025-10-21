package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	err           error
	flagNoColor   bool
	flagPageSize  int
	flagPageStyle string
	homeDir       string
	logger        *logrus.Logger
	rd            *raindrop.Raindrop
	rdbakConfig   string
	rdbakHome     string
	rdbakLogfile  string
	rootCmd       = &cobra.Command{
		Use:   "rdbak",
		Short: "rdbak is a command line utility to backup your raindrop.io bookmarks",
		Long:  "rdbak is a command line utility to backup your raindrop.io bookmarks",
	}
	validStyles = []string{"ascii", "light", "dark"}
	versionFull bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	homeDir, err = util.GetHome()
	if err != nil {
		fmt.Println("failed to determine home directory")
		os.Exit(1)
	}

	rdbakHome = filepath.Join(homeDir, ".config", "rdbak")

	err = util.VerifyDirectory(rdbakHome)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rdbakConfig = filepath.Join(rdbakHome, "config.yaml")
	rdbakLogfile = filepath.Join(rdbakHome, "rdbak.log")
	logger = util.ConfigureLogger(flagNoColor, rdbakLogfile)
}

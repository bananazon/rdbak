package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bananazon/raindrop/cmd/bookmarks"
	"github.com/bananazon/raindrop/cmd/collections"
	"github.com/bananazon/raindrop/cmd/encrypt"
	"github.com/bananazon/raindrop/cmd/tags"
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/bananazon/raindrop/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	err             error
	FlagNoColor     bool
	homeDir         string
	Logger          *logrus.Logger
	RD              *raindrop.Raindrop
	RaindropConfig  string
	RaindropHome    string
	RaindropLogFile string
	RootCmd         = &cobra.Command{
		Use:   "raindrop",
		Short: "Manage your raindrop.io bookmarks, collections, and tags",
		Long:  "Manage your raindrop.io bookmarks, collections, and tags",
	}

	versionFull bool
)

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	homeDir, err = util.GetHome()
	if err != nil {
		fmt.Println("failed to determine home directory")
		os.Exit(1)
	}

	RaindropHome = filepath.Join(homeDir, ".config", "raindrop")
	RaindropConfig = filepath.Join(RaindropHome, "config.yaml")

	err = util.VerifyDirectory(RaindropHome)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	RaindropLogFile = filepath.Join(RaindropHome, "raindrop.log")
	Logger = util.ConfigureLogger(FlagNoColor, RaindropLogFile)

	ctx := &context.AppContext{
		Logger:                    util.ConfigureLogger(FlagNoColor, RaindropLogFile),
		RaindropHome:              RaindropHome,
		RaindropConfig:            RaindropConfig,
		ValidCollectionsSortOrder: []string{"title", "-title", "-count"},
		ValidCollectionsViews:     []string{"grid", "list", "masonry", "simple"},
		ValidStyles:               []string{"ascii", "bright", "dark", "light"},
	}

	RootCmd.AddCommand(bookmarks.NewBookmarksCmd(ctx))
	RootCmd.AddCommand(collections.NewCollectionsCmd(ctx))
	RootCmd.AddCommand(encrypt.NewEncryptTokenCmd(ctx))
	RootCmd.AddCommand(tags.NewTagsCmd(ctx))
}

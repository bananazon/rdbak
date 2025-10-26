package bookmarks

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newBackupBookmarksCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "backup",
		Aliases: []string{"b"},
		Short:   "Back your raindrop.io bookmarks up to a YAML file",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			err := ctx.RD.BackupBookmarks(ctx.FlagPrune)
			if err != nil {
				ctx.Logger.Errorf("Failed to backup the bookmarks: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully backed up the bookmarks")
		},
	}

	ctx.GetBackupFlags(c)

	return c
}

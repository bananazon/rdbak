package bookmarks

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newRemoveBookmarkCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r"},
		Short:   "Remove an existing bookmark from your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			_, err := ctx.RD.API.RemoveBookmark(
				ctx.FlagRemoveBookmarkId,
			)
			if err != nil {
				ctx.Logger.Errorf("Failed to remove the bookmark: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully removed the bookmark")
		},
	}

	ctx.GetRemoveBookmarkFlags(c)

	return c
}

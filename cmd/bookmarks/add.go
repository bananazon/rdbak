package bookmarks

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newAddBookmarkCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new bookmark to your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			dataAddBookmarkResult, err := ctx.RD.API.AddBookmark(data.AddBookmarkPayload{
				CollectionId: int64(ctx.FlagAddBookmarkCollectionId),
				Excerpt:      ctx.FlagAddBookmarkExcerpt,
				Link:         ctx.FlagAddBookmarkLink,
				Tags:         ctx.FlagAddBookmarkTag,
				Title:        ctx.FlagAddBookmarkTitle,
			})
			if err != nil {
				ctx.Logger.Errorf("Failed to add the bookmark: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully added the bookmark with ID %d", dataAddBookmarkResult.Item.Id)
		},
	}

	ctx.GetAddBookmarkFlags(c)

	return c
}

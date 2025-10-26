package bookmarks

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newUpdateBookmarkCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "Update an existing bookmark in your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			_, err := ctx.RD.API.UpdateBookmark(
				ctx.FlagUpdateBookmarkBookmarkId,
				data.UpdateBookmarkPayload{
					CollectionId: int64(ctx.FlagUpdateBookmarkCollectionId),
					Excerpt:      ctx.FlagUpdateBookmarkExcerpt,
					Highlights:   ctx.FlagUpdateBookmarkHighlight,
					Important:    ctx.FlagUpdateBookmarkImportant,
					Link:         ctx.FlagUpdateBookmarkLink,
					Note:         ctx.FlagUpdateBookmarkNote,
					Tags:         ctx.FlagUpdateBookmarkTag,
					Title:        ctx.FlagUpdateBookmarkTitle,
				},
			)
			if err != nil {
				ctx.Logger.Errorf("Failed to update the bookmark: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully updated the bookmark")
		},
	}

	ctx.GetUpdateBookmarkFlags(c)

	return c
}

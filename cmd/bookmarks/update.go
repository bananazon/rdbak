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
		PreRunE: func(cmdC *cobra.Command, args []string) error {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Println("Failed to initialize raindrop:", err.Error())
				return err
			}
			ctx.RD = rd
			return nil
		},
		RunE: func(cmdC *cobra.Command, args []string) error {
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
				ctx.Logger.Println("Failed to add the bookmark:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully added the bookmark")
			return nil
		},
	}

	ctx.GetUpdateBookmarkFlags(c)

	return c
}

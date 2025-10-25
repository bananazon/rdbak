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
			_, err := ctx.RD.API.AddBookmark(data.AddBookmarkPayload{
				CollectionId: int64(ctx.FlagAddBookmarkCollectionId),
				Excerpt:      ctx.FlagAddBookmarkExcerpt,
				Link:         ctx.FlagAddBookmarkLink,
				Tags:         ctx.FlagAddBookmarkTag,
				Title:        ctx.FlagAddBookmarkTitle,
			})
			if err != nil {
				ctx.Logger.Println("Failed to add the bookmark:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully added the bookmark")
			return nil
		},
	}

	ctx.GetAddBookmarkFlags(c)

	return c
}

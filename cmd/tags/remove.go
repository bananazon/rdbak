package tags

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newRemoveTagsCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "remove",
		Short: "Remove tags from your raindrop.io account",
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
			_, err := ctx.RD.API.RemoveTags(data.RemoveTagsPayload{
				CollectionId: ctx.FlagRenameTagCollectionId,
				Tags:         ctx.FlagRemoveTagsTagNames,
			})
			if err != nil {
				ctx.Logger.Println("Failed to remove the tags:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully removed the tags")
			return nil
		},
	}

	ctx.GetRemoveTagsFlags(c)

	return c
}

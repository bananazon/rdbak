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
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			_, err := ctx.RD.API.RemoveTags(data.RemoveTagsPayload{
				CollectionId: ctx.FlagRenameTagCollectionId,
				Tags:         ctx.FlagRemoveTagsTagNames,
			})

			if err != nil {
				ctx.Logger.Errorf("Failed to remove the tags: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully removed the tags")
		},
	}

	ctx.GetRemoveTagsFlags(c)

	return c
}

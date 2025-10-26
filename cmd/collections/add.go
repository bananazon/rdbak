package collections

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newAddCollectionCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new collection to your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			addCollectionResult, err := ctx.RD.API.AddCollection(data.AddCollectionPayload{
				Title:  ctx.FlagAddCollectionTitle,
				Parent: ctx.FlagAddCollectionParent,
				Public: ctx.FlagAddCollectionPublic,
				View:   ctx.FlagAddCollectionView,
			})
			if err != nil {
				ctx.Logger.Errorf("Failed to add the collection: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully added the collection with ID %d", addCollectionResult.Item.Id)
		},
	}

	ctx.GetAddCollectionFlags(c)

	return c
}

package collections

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newUpdateCollectionCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "Update an existing collection in your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			_, err := ctx.RD.API.UpdateCollection(
				ctx.FlagUpdateCollectionCollectionId,
				data.UpdateCollectionPayload{
					Title:  ctx.FlagUpdateCollectionTitle,
					Parent: ctx.FlagUpdateCollectionParent,
					Public: ctx.FlagUpdateCollectionPublic,
					View:   ctx.FlagUpdateCollectionView,
				},
			)
			if err != nil {
				ctx.Logger.Errorf("Failed to update the collection: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully updated the collection")
		},
	}

	ctx.GetUpdateCollectionFlags(c)

	return c
}

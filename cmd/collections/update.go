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
				ctx.Logger.Println("Failed to update the collection:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully updated the collection")
			return nil
		},
	}

	ctx.GetUpdateCollectionFlags(c)

	return c
}

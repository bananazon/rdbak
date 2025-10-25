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
			_, err := ctx.RD.API.AddCollection(data.AddCollectionPayload{
				Title:  ctx.FlagAddCollectionTitle,
				Parent: ctx.FlagAddCollectionParent,
				Public: ctx.FlagAddCollectionPublic,
			})
			if err != nil {
				ctx.Logger.Println("Failed to add the collection:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully added the collection")
			return nil
		},
	}

	ctx.GetAddCollectionFlags(c)

	return c
}

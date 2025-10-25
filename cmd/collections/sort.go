package collections

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newSortCollectionsCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "sort",
		Aliases: []string{"s"},
		Short:   "Sort your raindrop.io collections",
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
			_, err := ctx.RD.API.SortCollections(data.SortCollectionPayload{
				Sort: ctx.FlagCollectionsSortOrder,
			})
			if err != nil {
				ctx.Logger.Println("Failed to sort the collection:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully sorted the collection")
			return nil
		},
	}

	ctx.GetSortCollectionFlags(c)

	return c
}

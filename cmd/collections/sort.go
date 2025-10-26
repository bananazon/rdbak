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
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			_, err := ctx.RD.API.SortCollections(data.SortCollectionPayload{
				Sort: ctx.FlagCollectionsSortOrder,
			})
			if err != nil {
				ctx.Logger.Errorf("Failed to sort the collections: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully sorted the collections")
		},
	}

	ctx.GetSortCollectionFlags(c)

	return c
}

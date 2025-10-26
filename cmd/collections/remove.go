package collections

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newRemoveCollectionCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r"},
		Short:   "Remove an existing collection from your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			_, err := ctx.RD.API.RemoveCollection(
				ctx.FlagRemoveCollectionId,
			)
			if err != nil {
				ctx.Logger.Errorf("Failed to remove the collection: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully removed the collection")
		},
	}

	ctx.GetRemoveCollectionFlags(c)

	return c
}

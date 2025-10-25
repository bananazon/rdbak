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
			_, err := ctx.RD.API.RemoveCollection(
				ctx.FlagRemoveCollectionId,
			)
			if err != nil {
				ctx.Logger.Println("Failed to remove the collection:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully removed the collection")
			return nil
		},
	}

	ctx.GetRemoveCollectionFlags(c)

	return c
}

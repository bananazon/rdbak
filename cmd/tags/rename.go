package tags

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newRenameTagCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "rename",
		Short: "Rename a tag in your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			_, err := ctx.RD.API.RenameTag(data.RenameTagPayload{
				CollectionId: ctx.FlagRenameTagCollectionId,
				OldName:      []string{ctx.FlagRenameTagOldName},
				NewName:      ctx.FlagRenameTagNewName,
			})
			if err != nil {
				ctx.Logger.Errorf("Failed to rename the tag: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.Logger.Infof("Successfully renamed the tag")
		},
	}

	ctx.GetRenameTagFlags(c)

	return c
}

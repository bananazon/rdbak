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
			_, err := ctx.RD.API.RenameTag(data.RenameTagPayload{
				CollectionId: ctx.FlagRenameTagCollectionId,
				OldName:      []string{ctx.FlagRenameTagOldName},
				NewName:      ctx.FlagRenameTagNewName,
			})
			if err != nil {
				ctx.Logger.Println("Failed to rename the tag:", err.Error())
				return err
			}
			ctx.Logger.Println("Successfully renamed the tag")
			return nil
		},
	}

	ctx.GetRenameTagFlags(c)

	return c
}

package collections

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newBackupCollectionsCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "backup",
		Aliases: []string{"b"},
		Short:   "Back your raindrop.io collections up to a YAML file",
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
			err := ctx.RD.BackupCollections(ctx.FlagPrune)
			if err != nil {
				ctx.Logger.Println("Backup failed:", err.Error())
				return err
			}
			ctx.Logger.Println("Backup successful")
			return nil
		},
	}

	ctx.GetBackupFlags(c)

	return c
}

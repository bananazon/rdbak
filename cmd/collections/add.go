package collections

import (
	"github.com/bananazon/rdbak/pkg/context"
	"github.com/bananazon/rdbak/pkg/raindrop"
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
				ctx.Logger.Println("Failed to initialize raindrop:", err)
				return err
			}
			ctx.RD = rd
			return nil
		},
		RunE: func(cmdC *cobra.Command, args []string) error {
			_, err := ctx.RD.API.AddCollection(
				ctx.FlagAddCollectionTitle,
				ctx.FlagAddCollectionParent,
				ctx.FlagAddCollectionPublic,
			)
			if err != nil {
				ctx.Logger.Println("AddCollection failed:", err)
				return err
			}
			ctx.Logger.Println("Collection added successfully.")
			return nil
		},
	}

	ctx.GetAddCollectionFlags(c)

	return c
}

// var (
// 	addCollectionCmd = &cobra.Command{
// 		Use:          "add",
// 		Aliases:      []string{"a"},
// 		Short:        "Add a new collection to your raindrop.io account",
// 		Long:         "Add a new collection to your raindrop.io account",
// 		PreRun:       addCollectionPreRunCmd,
// 		Run:          addCollectionRunCmd,
// 		SilenceUsage: false,
// 	}
// )

// func init() {
// 	cmd.GetAddCollectionFlags(addCollectionCmd)
// 	CollectionsCmd.AddCommand(addCollectionCmd)
// }

// func addCollectionPreRunCmd(cmdC *cobra.Command, args []string) {
// 	cmd.RD, err = raindrop.New(cmd.RaindropHome, cmd.RaindropConfig, cmd.Logger)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

// func addCollectionRunCmd(cmdC *cobra.Command, args []string) {
// 	_, err := cmd.RD.API.AddCollection(cmd.FlagAddCollectionTitle, cmd.FlagAddCollectionParent, cmd.FlagAddCollectionPublic)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

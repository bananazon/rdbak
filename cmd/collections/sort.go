package collections

import (
	"github.com/bananazon/rdbak/pkg/context"
	"github.com/bananazon/rdbak/pkg/raindrop"
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
				ctx.Logger.Println("Failed to initialize raindrop:", err)
				return err
			}
			ctx.RD = rd
			return nil
		},
		RunE: func(cmdC *cobra.Command, args []string) error {
			_, err := ctx.RD.API.SortCollections(ctx.FlagCollectionsSortOrder)
			if err != nil {
				ctx.Logger.Println("SortCollections failed:", err)
				return err
			}
			return nil
		},
	}

	ctx.GetSortCollectionFlags(c)

	return c
}

// var (
// 	sortCollectionCmd = &cobra.Command{
// 		Use:          "sort",
// 		Aliases:      []string{"s"},
// 		Short:        "Sort your raindrop.io collections",
// 		Long:         "Sort your raindrop.io collections",
// 		PreRun:       sortCollectionPreRunCmd,
// 		Run:          sortCollectionRunCmd,
// 		SilenceUsage: false,
// 	}
// )

// func init() {
// 	cmd.GetSortCollectionFlags(sortCollectionCmd)
// 	CollectionsCmd.AddCommand(sortCollectionCmd)
// }

// func sortCollectionPreRunCmd(cmdC *cobra.Command, args []string) {
// 	cmd.RD, err = raindrop.New(cmd.RaindropHome, cmd.RaindropConfig, cmd.Logger)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

// func sortCollectionRunCmd(cmdC *cobra.Command, args []string) {
// 	_, err = cmd.RD.API.SortCollections(cmd.FlagCollectionsSortOrder)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

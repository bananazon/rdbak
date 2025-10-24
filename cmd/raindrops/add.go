package raindrops

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newAddRaindropCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new raindrop to your raindrop.io account",
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
			_, err := ctx.RD.API.AddRaindrop(
				ctx.FlagAddRaindropLink,
				ctx.FlagAddRaindropTitle,
				int64(ctx.FlagAddRaindropCollectionId),
			)
			if err != nil {
				ctx.Logger.Println("AddRaindrop failed:", err)
				return err
			}
			ctx.Logger.Println("Raindrop added successfully.")
			return nil
		},
	}

	ctx.GetAddRaindropFlags(c)

	return c
}

// var (
// 	addRaindropCmd = &cobra.Command{
// 		Use:          "add",
// 		Aliases:      []string{"a"},
// 		Short:        "Add a new raindrop to your raindrop.io account",
// 		Long:         "Add a new raindrop to your raindrop.io account",
// 		PreRun:       addRaindropPreRunCmd,
// 		Run:          addRaindropRunCmd,
// 		SilenceUsage: false,
// 	}
// )

// func init() {
// 	cmd.GetAddRaindropFlags(addRaindropCmd)
// 	raindropCmd.AddCommand(addRaindropCmd)
// }

// func addRaindropPreRunCmd(cmdC *cobra.Command, args []string) {
// 	cmd.RD, err = raindrop.New(cmd.RaindropHome, cmd.RaindropConfig, cmd.Logger)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

// func addRaindropRunCmd(cmdC *cobra.Command, args []string) {
// 	_, err := cmd.RD.API.AddRaindrop(cmd.FlagAddRaindropLink, cmd.FlagAddRaindropTitle, int64(cmd.FlagAddRaindropCollectionId))
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

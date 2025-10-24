package raindrops

import (
	"github.com/bananazon/rdbak/pkg/context"
	"github.com/bananazon/rdbak/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newBackupRaindropsCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "backup",
		Aliases: []string{"b"},
		Short:   "Back your raindrop.io bookmarks up to a YAML file",
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
			err := ctx.RD.Backup(ctx.FlagPrune)
			if err != nil {
				ctx.Logger.Println("Backup failed:", err)
				return err
			}
			ctx.Logger.Println("Backup successful.")
			return nil
		},
	}

	ctx.GetBackupRaindropsFlags(c)

	return c
}

// var (
// 	backupCmd = &cobra.Command{
// 		Use:          "backup",
// 		Aliases:      []string{"b"},
// 		Short:        "Back your raindrop.io bookmarks up to a YAML file",
// 		Long:         "Back your raindrop.io bookmarks up to a YAML file",
// 		PreRun:       backupPreRunCmd,
// 		Run:          backupRunCmd,
// 		SilenceUsage: false,
// 	}
// )

// func init() {
// 	cmd.GetBackupFlags(backupCmd)
// 	raindropCmd.AddCommand(backupCmd)
// }

// func backupPreRunCmd(cmdC *cobra.Command, args []string) {
// 	cmd.RD, err = raindrop.New(cmd.RaindropHome, cmd.RaindropConfig, cmd.Logger)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

// func backupRunCmd(cmdC *cobra.Command, args []string) {
// 	err = cmd.RD.Backup(cmd.FlagPrune)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

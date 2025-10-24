package bookmarks

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

func newAddBookmarkCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a new bookmark to your raindrop.io account",
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
			_, err := ctx.RD.API.AddBookmark(
				ctx.FlagAddBookmarkLink,
				ctx.FlagAddBookmarkTitle,
				int64(ctx.FlagAddBookmarkCollectionId),
			)
			if err != nil {
				ctx.Logger.Println("Failed to add the bookmark:", err)
				return err
			}
			ctx.Logger.Println("Bookmark added successfully.")
			return nil
		},
	}

	ctx.GetAddBookmarkFlags(c)

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
// 	_, err := cmd.RD.API.AddRaindrop(cmd.FlagAddBookmarkLink, cmd.FlagAddBookmarkTitle, int64(cmd.FlagAddBookmarkCollectionId))
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

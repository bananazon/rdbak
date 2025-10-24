package raindrops

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/spf13/cobra"
)

func NewRaindropsCmd(ctx *context.AppContext) (cmdC *cobra.Command) {
	cmdC = &cobra.Command{
		Use:     "raindrops",
		Aliases: []string{"r"},
		Short:   "Manage bookmarks in your raindrop.io account",
		Long:    "Manage bookmarks in your raindrop.io account",
	}

	cmdC.AddCommand(newAddRaindropCmd(ctx))
	cmdC.AddCommand(newBackupRaindropsCmd(ctx))
	cmdC.AddCommand(newListRaindropsCmd(ctx))
	cmdC.AddCommand(newRemoveRaindropCmd(ctx))

	return cmdC
}

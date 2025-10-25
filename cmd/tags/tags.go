package tags

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/spf13/cobra"
)

func NewTagsCmd(ctx *context.AppContext) (cmdC *cobra.Command) {
	cmdC = &cobra.Command{
		Use:     "tags",
		Aliases: []string{"t"},
		Short:   "Manage tags in your raindrop.io account",
		Long:    "Manage tags in your raindrop.io account",
	}

	cmdC.AddCommand(newListTagsCmd(ctx))
	cmdC.AddCommand(newRemoveTagsCmd(ctx))
	cmdC.AddCommand(newRenameTagCmd(ctx))

	return cmdC
}

package bookmarks

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/spf13/cobra"
)

func NewBookmarksCmd(ctx *context.AppContext) (cmdC *cobra.Command) {
	cmdC = &cobra.Command{
		Use:     "bookmarks",
		Aliases: []string{"b"},
		Short:   "Manage bookmarks in your raindrop.io account",
		Long:    "Manage bookmarks in your raindrop.io account",
	}

	cmdC.AddCommand(newAddBookmarkCmd(ctx))
	cmdC.AddCommand(newBackupBookmarksCmd(ctx))
	cmdC.AddCommand(newListBookmarksCmd(ctx))
	cmdC.AddCommand(newRemoveBookmarkCmd(ctx))

	return cmdC
}

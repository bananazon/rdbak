package collections

import (
	"github.com/bananazon/raindrop/pkg/context"
	"github.com/spf13/cobra"
)

func NewCollectionsCmd(ctx *context.AppContext) (cmdC *cobra.Command) {
	cmdC = &cobra.Command{
		Use:     "collections",
		Aliases: []string{"c"},
		Short:   "Manage collections in your raindrop.io account",
		Long:    "Manage collections in your raindrop.io account",
	}

	cmdC.AddCommand(newAddCollectionCmd(ctx))
	cmdC.AddCommand(newBackupCollectionsCmd(ctx))
	cmdC.AddCommand(newListCollectionsCmd(ctx))
	cmdC.AddCommand(newRemoveCollectionCmd(ctx))
	cmdC.AddCommand(newSortCollectionsCmd(ctx))
	cmdC.AddCommand(newUpdateCollectionCmd(ctx))

	return cmdC
}

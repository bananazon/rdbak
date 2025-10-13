package cmd

import (
	"github.com/spf13/cobra"
)

func GetBackupFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&flagPrune, "prune", "p", false, "Prune older bookmarks-{timestamp}.yaml files")
}

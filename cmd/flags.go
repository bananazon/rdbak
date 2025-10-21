package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func GetBackupFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&flagPrune, "prune", "p", false, "Prune older bookmarks-{timestamp}.yaml files")
}

func GetTableFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&flagPageStyle, "page-style", "s", "ascii", fmt.Sprintf("The page style to use; one of %s", strings.Join(validStyles, ",")))
	cmd.Flags().IntVarP(&flagPageSize, "page-size", "p", 40, "The page size for the paginator")
}

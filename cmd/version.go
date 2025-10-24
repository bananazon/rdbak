package cmd

import (
	"fmt"

	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:          "version",
		Short:        "Print the current raindrop version",
		Long:         "Print the current raindrop version",
		RunE:         runVersionCmd,
		SilenceUsage: true,
	}
)

func init() {
	versionCmd.Flags().BoolVarP(&versionFull, "full", "f", false, "Display more version information")
	RootCmd.AddCommand(versionCmd)
}

func runVersionCmd(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", raindrop.Version(true, true, versionFull))

	return nil
}

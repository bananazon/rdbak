package cmd

import (
	"os"
	"strconv"

	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"l", "ls"},
		Short:        "List the bookmarks in your raindrop.io account",
		Long:         "List the bookmarks in your raindrop.io account",
		PreRun:       listPreRunCmd,
		Run:          listRunCmd,
		SilenceUsage: false,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func listPreRunCmd(cmd *cobra.Command, args []string) {
	rd, err = raindrop.New(rdbakHome, rdbakConfig, flagPrune, logger)
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

func listRunCmd(cmd *cobra.Command, args []string) {
	raindrops, err := rd.List()
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}

	tableOut := [][]string{}
	for _, raindrop := range raindrops {
		tableOut = append(tableOut, []string{
			strconv.Itoa(int(raindrop.Id)),
			raindrop.Link,
		})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Link"})
	table.Bulk(tableOut)
	table.Render()
}

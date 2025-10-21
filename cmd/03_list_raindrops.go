package cmd

import (
	"os"
	"strconv"
	"strings"

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
			raindrop.Type,
			raindrop.Link,
			strings.Join(raindrop.Tags, ","),
		})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Type", "Link", "Tags"})
	table.Bulk(tableOut)
	table.Render()
}

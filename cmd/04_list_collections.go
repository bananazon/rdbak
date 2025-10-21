package cmd

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	listCollectionsCmd = &cobra.Command{
		Use:          "list-collections",
		Aliases:      []string{"lc"},
		Short:        "List the collections in your raindrop.io account",
		Long:         "List the collections in your raindrop.io account",
		PreRun:       listCollectionsPreRunCmd,
		Run:          listCollectionsRunCmd,
		SilenceUsage: false,
	}
)

func init() {
	rootCmd.AddCommand(listCollectionsCmd)
}

func listCollectionsPreRunCmd(cmd *cobra.Command, args []string) {

}

func listCollectionsRunCmd(cmd *cobra.Command, args []string) {
	tableOut := [][]string{}
	for _, collection := range rd.Collections {
		var description string

		if collection.Description == "" {
			description = "N/A"
		} else {
			description = collection.Description
		}

		tableOut = append(tableOut, []string{
			strconv.Itoa(int(collection.Id)),
			collection.Title,
			collection.View,
			description,
		})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Title", "View", "Description"})
	table.Bulk(tableOut)
	table.Render()
}

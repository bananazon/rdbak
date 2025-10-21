package cmd

import (
	"fmt"
	"strconv"

	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
	GetTableFlags(listCollectionsCmd)
	rootCmd.AddCommand(listCollectionsCmd)
}

func listCollectionsPreRunCmd(cmd *cobra.Command, args []string) {
	rd, err = raindrop.New(rdbakHome, rdbakConfig, logger)
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

func listCollectionsRunCmd(cmd *cobra.Command, args []string) {
	collections, err := rd.ListCollections()
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}

	t := table.NewWriter()
	switch flagPageStyle {
	case "bright":
		t.SetStyle(table.StyleColoredBright)
	case "dark":
		t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	case "light":
		t.SetStyle(table.StyleLight)
	}

	// t.SetColumnConfigs([]table.ColumnConfig{{Name: "Description", WidthMax: 80}})
	t.SortBy([]table.SortBy{
		{Name: "Title", Mode: table.Asc},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Count", Align: text.AlignLeft},
	})
	t.SetPageSize(flagPageSize)
	t.AppendHeader(table.Row{"ID", "Title", "View", "Description", "Count"})

	for _, collection := range collections {
		var description string

		if collection.Description == "" {
			description = "N/A"
		} else {
			description = collection.Description
		}
		t.AppendRow(table.Row{
			strconv.Itoa(int(collection.Id)),
			collection.Title,
			collection.View,
			description,
			collection.Count,
		})
	}

	fmt.Println(t.Render())
}

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	listRaindropsCmd = &cobra.Command{
		Use:          "list-raindrops",
		Aliases:      []string{"lr"},
		Short:        "List the raindrops in your raindrop.io account",
		Long:         "List the raindrops in your raindrop.io account",
		PreRun:       listRaindropsPreRunCmd,
		Run:          listRaindropsRunCmd,
		SilenceUsage: false,
	}
)

func init() {
	GetTableFlags(listRaindropsCmd)
	rootCmd.AddCommand(listRaindropsCmd)
}

func listRaindropsPreRunCmd(cmd *cobra.Command, args []string) {
	rd, err = raindrop.New(rdbakHome, rdbakConfig, logger)
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

func listRaindropsRunCmd(cmd *cobra.Command, args []string) {
	raindrops, err := rd.ListRaindrops()
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}

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

	// t.SetColumnConfigs([]table.ColumnConfig{{Name: "Link", WidthMax: 80}})
	t.SortBy([]table.SortBy{
		{Name: "Collection", Mode: table.Asc},
		{Name: "Link", Mode: table.Asc},
	})
	t.SetPageSize(flagPageSize)
	t.AppendHeader(table.Row{"ID", "Collection", "Link", "Tags"})

	for idx, raindrop := range raindrops {
		collectionId := raindrops[idx].Collection.Id
		collection, exists := collections[uint64(collectionId)]
		if exists {
			if raindrops[idx].Collection.Id > 0 {
				raindrops[idx].Collection.Name = collection.Title
			} else {
				raindrops[idx].Collection.Name = "N/A"
			}
		}
		t.AppendRow(table.Row{
			strconv.Itoa(int(raindrop.Id)),
			raindrop.Collection.Name,
			raindrop.Link,
			strings.Join(raindrop.Tags, ","),
		})
	}

	fmt.Println(t.Render())
}

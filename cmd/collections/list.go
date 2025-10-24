package collections

import (
	"fmt"
	"strconv"

	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	rdtable "github.com/bananazon/raindrop/pkg/table"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func newListCollectionsCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List the collections in your raindrop.io account",
		PreRunE: func(cmdC *cobra.Command, args []string) error {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Println("Failed to initialize raindrop:", err)
				return err
			}
			ctx.RD = rd
			return nil
		},
		RunE: func(cmdC *cobra.Command, args []string) error {
			collections, err := ctx.RD.ListCollections()
			if err != nil {
				ctx.Logger.Println("ListCollections failed:", err)
				return err
			}

			t := rdtable.GetTableTemplate("Collections", ctx.FlagPageSize, ctx.FlagPageStyle)
			t.SortBy([]table.SortBy{{Name: "Title", Mode: table.Asc}})
			t.SetColumnConfigs([]table.ColumnConfig{{Name: "Count", Align: text.AlignLeft}})
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

			return nil
		},
	}

	ctx.GetTableFlags(c)

	return c

}

// var (
// 	listCollectionsCmd = &cobra.Command{
// 		Use:          "list",
// 		Aliases:      []string{"l", "ls"},
// 		Short:        "List the collections in your raindrop.io account",
// 		Long:         "List the collections in your raindrop.io account",
// 		PreRun:       listCollectionsPreRunCmd,
// 		Run:          listCollectionsRunCmd,
// 		SilenceUsage: false,
// 	}
// )

// func init() {
// 	cmd.GetTableFlags(listCollectionsCmd)
// 	CollectionsCmd.AddCommand(listCollectionsCmd)
// }

// func listCollectionsPreRunCmd(cmdC *cobra.Command, args []string) {
// 	cmd.RD, err = raindrop.New(cmd.RaindropHome, cmd.RaindropConfig, cmd.Logger)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

// func listCollectionsRunCmd(cmdC *cobra.Command, args []string) {
// 	collections, err := cmd.RD.ListCollections()
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}

// 	t := cmd.GetTableTemplate("Collections")
// 	t.SortBy([]table.SortBy{{Name: "Title", Mode: table.Asc}})
// 	t.SetColumnConfigs([]table.ColumnConfig{{Name: "Count", Align: text.AlignLeft}})
// 	t.AppendHeader(table.Row{"ID", "Title", "View", "Description", "Count"})

// 	for _, collection := range collections {
// 		var description string

// 		if collection.Description == "" {
// 			description = "N/A"
// 		} else {
// 			description = collection.Description
// 		}
// 		t.AppendRow(table.Row{
// 			strconv.Itoa(int(collection.Id)),
// 			collection.Title,
// 			collection.View,
// 			description,
// 			collection.Count,
// 		})
// 	}

// 	fmt.Println(t.Render())
// }

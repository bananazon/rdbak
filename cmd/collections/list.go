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
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			collections, err := ctx.RD.ListCollections()
			if err != nil {
				ctx.Logger.Errorf("Failed to get a list of collections: %s", err.Error())
				ctx.Logger.Exit(1)
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
		},
	}

	ctx.GetTableFlags(c)

	return c

}

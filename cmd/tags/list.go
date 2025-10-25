package tags

import (
	"fmt"

	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	rdtable "github.com/bananazon/raindrop/pkg/table"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func newListTagsCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List the tags in your raindrop.io account",
		PreRunE: func(cmdC *cobra.Command, args []string) error {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Println("Failed to initialize raindrop:", err.Error())
				return err
			}
			ctx.RD = rd
			return nil
		},
		RunE: func(cmdC *cobra.Command, args []string) error {
			listTagsResult, err := ctx.RD.API.ListTags()
			if err != nil {
				ctx.Logger.Println("Failed to get a list of tags:", err.Error())
				return err
			}

			t := rdtable.GetTableTemplate("Tags", ctx.FlagPageSize, ctx.FlagPageStyle)
			t.SortBy([]table.SortBy{{Name: "ID", Mode: table.Asc}})
			t.SetColumnConfigs([]table.ColumnConfig{{Name: "Count", Align: text.AlignLeft}})
			t.AppendHeader(table.Row{"ID", "Count"})

			for _, tag := range listTagsResult.Items {
				t.AppendRow(table.Row{
					tag.Id,
					tag.Count,
				})
			}

			fmt.Println(t.Render())

			return nil
		},
	}

	ctx.GetTableFlags(c)

	return c

}

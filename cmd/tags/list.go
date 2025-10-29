package tags

import (
	"fmt"
	"os"

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
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			listTagsResult, err := ctx.RD.API.ListTags()
			if err != nil {
				ctx.Logger.Errorf("Failed to get a list of tags: %s", err.Error())
				ctx.Logger.Exit(1)
			}

			if ctx.FlagPageStyle == "list" {
				for _, tag := range listTagsResult.Items {
					fmt.Fprintf(os.Stdout, "%s = %s\n", "      id", tag.Id)
					fmt.Fprintf(os.Stdout, "%s = %d\n", "   count", tag.Count)
					fmt.Fprintln(os.Stdout, "")
				}
			} else {
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

				fmt.Fprintln(os.Stdout, t.Render())
			}
		},
	}

	ctx.GetTableFlags(c)

	return c

}

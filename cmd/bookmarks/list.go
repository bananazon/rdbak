package bookmarks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	rdtable "github.com/bananazon/raindrop/pkg/table"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func newListBookmarksCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List the bookmarks in your raindrop.io account",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			if ctx.ScreenWidth < 81 {
				ctx.Logger.Errorf("Screen width is only %d, please increase it and try again", ctx.ScreenWidth)
				ctx.Logger.Exit(1)
			}
			bookmarks, err := ctx.RD.ListBookmarks()
			if err != nil {
				ctx.Logger.Errorf("Failed to get a list of bookmarks: %s", err.Error())
				ctx.Logger.Exit(1)
			}

			collections, err := ctx.RD.ListCollections()
			if err != nil {
				ctx.Logger.Errorf("Failed to get a list of collections: %s", err.Error())
				ctx.Logger.Exit(1)
			}

			t := rdtable.GetTableTemplate("Bookmarks", ctx.FlagPageSize, ctx.FlagPageStyle)
			// t.SetColumnConfigs([]table.ColumnConfig{
			// 	{Name: "Link", WidthMax: 50},
			// })
			t.SortBy([]table.SortBy{{Name: "Collection", Mode: table.Asc}, {Name: "Link", Mode: table.Asc}})
			t.AppendHeader(table.Row{"ID", "Collection", "Link", "Tags"})

			for idx, raindrop := range bookmarks {
				bookmarks[idx].Collection.Name = "Unsorted"
				collectionId := bookmarks[idx].Collection.Id
				collection, exists := collections[uint64(collectionId)]
				if exists {
					if bookmarks[idx].Collection.Id > 0 {
						bookmarks[idx].Collection.Name = collection.Title
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
		},
	}

	ctx.GetTableFlags(c)

	return c
}

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
			bookmarks, err := ctx.RD.ListBookmarks()
			if err != nil {
				ctx.Logger.Println("Failed to get a list of bookmarks:", err.Error())
				return err
			}

			collections, err := ctx.RD.ListCollections()
			if err != nil {
				ctx.Logger.Println("Failed to get a list of collections:", err.Error())
				return err
			}

			t := rdtable.GetTableTemplate("Bookmarks", ctx.FlagPageSize, ctx.FlagPageStyle)
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

			return nil
		},
	}

	ctx.GetTableFlags(c)

	return c
}

package bookmarks

import (
	"fmt"
	"os"
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
			var useList = false
			if ctx.ScreenWidth < 80 {
				ctx.Logger.Infof("Screen width is only %d, using list output mode", ctx.ScreenWidth)
				useList = true
			}

			if ctx.FlagPageStyle == "list" {
				useList = true
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

			// Find and set collection name
			for idx, _ := range bookmarks {
				bookmarks[idx].Collection.Name = "Unsorted"
				collectionId := bookmarks[idx].Collection.Id
				collection, exists := collections[uint64(collectionId)]
				if exists {
					if bookmarks[idx].Collection.Id > 0 {
						bookmarks[idx].Collection.Name = collection.Title
					}
				}
			}

			if useList {
				for _, bookmark := range bookmarks {
					fmt.Fprintf(os.Stdout, "%s = %d\n", "           id", bookmark.Id)
					fmt.Fprintf(os.Stdout, "%s = %s\n", "   collection", bookmark.Collection.Name)
					fmt.Fprintf(os.Stdout, "%s = %s\n", "         link", bookmark.Link)
					fmt.Fprintf(os.Stdout, "%s = %s\n", "         tags", strings.Join(bookmark.Tags, ","))
					fmt.Fprintln(os.Stdout, "")
				}
			} else {
				maxIdWidth := 12
				maxCollectionWidth := 20
				maxTagsWidth := 30

				t := rdtable.GetTableTemplate("Bookmarks", ctx.FlagPageSize, ctx.FlagPageStyle)

				t.SetColumnConfigs([]table.ColumnConfig{
					{Name: "ID", WidthMax: maxIdWidth},
					{Name: "Collection", WidthMax: maxCollectionWidth},
					{Name: "Tags", WidthMax: maxTagsWidth},
					{Name: "Link", WidthMax: ctx.ScreenWidth - (maxIdWidth + maxCollectionWidth + maxTagsWidth)},
				})
				t.SortBy([]table.SortBy{{Name: "Collection", Mode: table.Asc}, {Name: "Link", Mode: table.Asc}})
				t.AppendHeader(table.Row{"ID", "Collection", "Link", "Tags"})

				for _, bookmark := range bookmarks {
					t.AppendRow(table.Row{
						strconv.Itoa(int(bookmark.Id)),
						bookmark.Collection.Name,
						bookmark.Link,
						strings.Join(bookmark.Tags, ","),
					})
				}

				fmt.Fprintln(os.Stdout, t.Render())
			}
		},
	}

	ctx.GetTableFlags(c)

	return c
}

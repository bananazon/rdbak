package raindrops

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

func newListRaindropsCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List the raindrops in your raindrop.io account",
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
			raindrops, err := ctx.RD.ListRaindrops()
			if err != nil {
				ctx.Logger.Error(err)
				ctx.Logger.Exit(1)
			}

			collections, err := ctx.RD.ListCollections()
			if err != nil {
				ctx.Logger.Error(err)
				ctx.Logger.Exit(1)
			}

			t := rdtable.GetTableTemplate("Raindrops", ctx.FlagPageSize, ctx.FlagPageStyle)
			t.SortBy([]table.SortBy{{Name: "Collection", Mode: table.Asc}, {Name: "Link", Mode: table.Asc}})
			t.AppendHeader(table.Row{"ID", "Collection", "Link", "Tags"})

			// pretty.Println(raindrops)
			// os.Exit(0)

			for idx, raindrop := range raindrops {
				raindrops[idx].Collection.Name = "Unsorted"
				collectionId := raindrops[idx].Collection.Id
				collection, exists := collections[uint64(collectionId)]
				if exists {
					if raindrops[idx].Collection.Id > 0 {
						raindrops[idx].Collection.Name = collection.Title
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

// var (
// 	listRaindropsCmd = &cobra.Command{
// 		Use:          "list",
// 		Aliases:      []string{"l", "ls"},
// 		Short:        "List the raindrops in your raindrop.io account",
// 		Long:         "List the raindrops in your raindrop.io account",
// 		PreRun:       listRaindropsPreRunCmd,
// 		Run:          listRaindropsRunCmd,
// 		SilenceUsage: false,
// 	}
// )

// func init() {
// 	cmd.GetTableFlags(listRaindropsCmd)
// 	raindropCmd.AddCommand(listRaindropsCmd)
// }

// func listRaindropsPreRunCmd(cmdC *cobra.Command, args []string) {
// 	cmd.RD, err = raindrop.New(cmd.RaindropHome, cmd.RaindropConfig, cmd.Logger)
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}
// }

// func listRaindropsRunCmd(cmdC *cobra.Command, args []string) {
// 	raindrops, err := cmd.RD.ListRaindrops()
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}

// 	collections, err := cmd.RD.ListCollections()
// 	if err != nil {
// 		cmd.Logger.Error(err)
// 		cmd.Logger.Exit(1)
// 	}

// 	t := cmd.GetTableTemplate("Raindrops")
// 	t.SortBy([]table.SortBy{{Name: "Collection", Mode: table.Asc}, {Name: "Link", Mode: table.Asc}})
// 	t.AppendHeader(table.Row{"ID", "Collection", "Link", "Tags"})

// 	for idx, raindrop := range raindrops {
// 		collectionId := raindrops[idx].Collection.Id
// 		collection, exists := collections[uint64(collectionId)]
// 		if exists {
// 			if raindrops[idx].Collection.Id > 0 {
// 				raindrops[idx].Collection.Name = collection.Title
// 			} else {
// 				raindrops[idx].Collection.Name = "N/A"
// 			}
// 		}
// 		t.AppendRow(table.Row{
// 			strconv.Itoa(int(raindrop.Id)),
// 			raindrop.Collection.Name,
// 			raindrop.Link,
// 			strings.Join(raindrop.Tags, ","),
// 		})
// 	}

// 	fmt.Println(t.Render())
// }

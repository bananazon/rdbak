package context

import (
	"fmt"
	"strings"

	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type AppContext struct {
	FlagAddCollectionParent     int64
	FlagAddCollectionPublic     bool
	FlagAddCollectionTitle      string
	FlagAddBookmarkCollectionId int64
	FlagAddBookmarkExcerpt      string
	FlagAddBookmarkLink         string
	FlagAddBookmarkTag          []string
	FlagAddBookmarkTitle        string
	FlagCollectionsSortOrder    string
	FlagNoColor                 bool
	FlagPageSize                int
	FlagPageStyle               string
	FlagPrune                   bool
	FlagRemoveCollectionId      int64
	FlagRemoveBookmarkId        int64
	Logger                      *logrus.Logger
	RaindropConfig              string
	RaindropHome                string
	RD                          *raindrop.Raindrop
	ValidCollectionsSortOrder   []string
	ValidStyles                 []string
}

func (ac *AppContext) GetBackupFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&ac.FlagPrune, "prune", "p", false, "Prune older [bookmarks|collections]-{timestamp}.yaml files")
}

func (ac *AppContext) GetTableFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ac.FlagPageStyle, "page-style", "s", "light", fmt.Sprintf("The page style to use; one of %s", strings.Join(ac.ValidStyles, ",")))
	cmd.Flags().IntVarP(&ac.FlagPageSize, "page-size", "p", 40, "The page size for the paginator")
}

func (ac *AppContext) GetAddCollectionFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ac.FlagAddCollectionTitle, "title", "t", "", "The title for the new collection")
	cmd.Flags().Int64VarP(&ac.FlagAddCollectionParent, "parent", "p", 0, "The parent ID of the new collection")
	cmd.Flags().BoolVarP(&ac.FlagAddCollectionPublic, "public", "", false, "Set the new collection to private")

	cmd.MarkFlagRequired("title")
}

func (ac *AppContext) GetRemoveCollectionFlags(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&ac.FlagRemoveCollectionId, "id", "i", -1, "ID of the collecion to remove")

	cmd.MarkFlagRequired("id")
}

func (ac *AppContext) GetSortCollectionFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ac.FlagCollectionsSortOrder, "order", "o", "title", fmt.Sprintf("Sort the collections; one of %s", strings.Join(ac.ValidCollectionsSortOrder, ",")))
}

func (ac *AppContext) GetAddBookmarkFlags(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&ac.FlagAddBookmarkCollectionId, "collection", "c", -1, "The collection ID of the new raindrop (`raindrop lc` to find the ID)")
	// cmd.Flags().StringVarP(&ac.FlagAddBookmarkExcerpt, "excerpt", "e", "", "A brief description of the link")
	cmd.Flags().StringVarP(&ac.FlagAddBookmarkTitle, "title", "t", "", "The title for the new raindrop")
	cmd.Flags().StringVarP(&ac.FlagAddBookmarkLink, "link", "l", "", "The URL for the new raindrop")
	cmd.Flags().StringSliceVarP(&ac.FlagAddBookmarkTag, "tag", "", []string{}, "One or more tags to use")

	cmd.MarkFlagRequired("link")
	cmd.MarkFlagRequired("title")
}

func (ac *AppContext) GetRemoveBookmarkFlags(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&ac.FlagRemoveBookmarkId, "id", "i", -1, "ID of the bookmark to remove")

	cmd.MarkFlagRequired("id")
}

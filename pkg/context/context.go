package context

import (
	"fmt"
	"strings"

	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type AppContext struct {
	FlagAddBookmarkCollectionId    int64
	FlagAddBookmarkExcerpt         string
	FlagAddBookmarkHighlight       []string
	FlagAddBookmarkImportant       bool
	FlagAddBookmarkLink            string
	FlagAddBookmarkNote            string
	FlagAddBookmarkTag             []string
	FlagAddBookmarkTitle           string
	FlagAddCollectionParent        int64
	FlagAddCollectionPublic        bool
	FlagAddCollectionTitle         string
	FlagCollectionsSortOrder       string
	FlagNoColor                    bool
	FlagPageSize                   int
	FlagPageStyle                  string
	FlagPrune                      bool
	FlagRemoveBookmarkId           int64
	FlagRemoveCollectionId         int64
	FlagRemoveTagsCollectionId     int64
	FlagRemoveTagsTagNames         []string
	FlagRenameTagCollectionId      int64
	FlagRenameTagNewName           string
	FlagRenameTagOldName           string
	FlagUpdateBookmarkBookmarkId   int64
	FlagUpdateBookmarkCollectionId int64
	FlagUpdateBookmarkExcerpt      string
	FlagUpdateBookmarkHighlight    []string
	FlagUpdateBookmarkImportant    bool
	FlagUpdateBookmarkLink         string
	FlagUpdateBookmarkNote         string
	FlagUpdateBookmarkTag          []string
	FlagUpdateBookmarkTitle        string
	Logger                         *logrus.Logger
	RaindropConfig                 string
	RaindropHome                   string
	RD                             *raindrop.Raindrop
	ValidCollectionsSortOrder      []string
	ValidStyles                    []string
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
	cmd.Flags().Int64VarP(&ac.FlagAddBookmarkCollectionId, "collection", "c", -1, "The collection ID of the new raindrop (use 'raindrop collections list' to find the ID)")
	cmd.Flags().StringVarP(&ac.FlagAddBookmarkExcerpt, "excerpt", "e", "", "A brief description of the link")
	cmd.Flags().StringVarP(&ac.FlagAddBookmarkTitle, "title", "t", "", "The title for the new raindrop")
	cmd.Flags().StringVarP(&ac.FlagAddBookmarkLink, "link", "l", "", "The URL for the new raindrop")
	cmd.Flags().StringVarP(&ac.FlagAddBookmarkNote, "note", "n", "", "Add a note to the bookmark")
	cmd.Flags().StringSliceVar(&ac.FlagAddBookmarkHighlight, "highlight", []string{}, "Bookmark highlight to set; can be used more than once")
	cmd.Flags().BoolVarP(&ac.FlagAddBookmarkImportant, "important", "i", false, "Set the important flag to true")
	cmd.Flags().StringSliceVar(&ac.FlagAddBookmarkTag, "tag", []string{}, "The tag to use; can be used more than once")

	cmd.MarkFlagRequired("link")
	cmd.MarkFlagRequired("title")
}

func (ac *AppContext) GetRemoveBookmarkFlags(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&ac.FlagRemoveBookmarkId, "id", "i", -1, "ID of the bookmark to remove")

	cmd.MarkFlagRequired("id")
}

func (ac *AppContext) GetUpdateBookmarkFlags(cmd *cobra.Command) {
	cmd.Flags().Int64Var(&ac.FlagUpdateBookmarkBookmarkId, "id", -1, "The bookmark ID (use 'raindrop bookmarks list' to find the ID)")
	cmd.Flags().Int64VarP(&ac.FlagUpdateBookmarkCollectionId, "collection", "c", -1, "The collection ID of the new raindrop (use 'raindrop collections list' to find the ID)")
	cmd.Flags().StringVarP(&ac.FlagUpdateBookmarkExcerpt, "excerpt", "e", "", "A brief description of the link")
	cmd.Flags().StringVarP(&ac.FlagUpdateBookmarkTitle, "title", "t", "", "The title for the new raindrop")
	cmd.Flags().StringVarP(&ac.FlagUpdateBookmarkLink, "link", "l", "", "The URL for the new raindrop")
	cmd.Flags().StringVarP(&ac.FlagUpdateBookmarkNote, "note", "n", "", "Add a note to the bookmark")
	cmd.Flags().StringSliceVar(&ac.FlagUpdateBookmarkHighlight, "highlight", []string{}, "Bookmark highlight to set; can be used more than once")
	cmd.Flags().BoolVarP(&ac.FlagUpdateBookmarkImportant, "important", "i", false, "Set the important flag to true")
	cmd.Flags().StringSliceVar(&ac.FlagUpdateBookmarkTag, "tag", []string{}, "The tag to use; can be used more than once")

	cmd.MarkFlagRequired("id")
}

func (ac *AppContext) GetRemoveTagsFlags(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&ac.FlagRemoveTagsCollectionId, "collection", "c", -1, "Restrict tag replacement to a specific collection ID (use 'raindrop collections list' to find the ID)")
	cmd.Flags().StringVarP(&ac.FlagRenameTagOldName, "old-name", "o", "", "The existing tag name")
	cmd.Flags().StringVarP(&ac.FlagRenameTagNewName, "new-name", "n", "", "The new tag name")

	cmd.MarkFlagRequired("old-name")
	cmd.MarkFlagRequired("new-name")
}

func (ac *AppContext) GetRenameTagFlags(cmd *cobra.Command) {
	cmd.Flags().Int64VarP(&ac.FlagRenameTagCollectionId, "collection", "c", -1, "Restrict tag replacement to a specific collection ID (use 'raindrop collections list' to find the ID)")
	cmd.Flags().StringSliceVarP(&ac.FlagRemoveTagsTagNames, "tag", "t", []string{}, "Tag name to remove; can be used more than once")

	cmd.MarkFlagRequired("tag")
}

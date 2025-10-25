package raindrop

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/bananazon/raindrop/pkg/api"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/bananazon/raindrop/pkg/util"
	"gopkg.in/yaml.v3"
)

func (r *Raindrop) BackupBookmarks(flagPrune bool) (err error) {
	r.PruneOlder = flagPrune

	r.Logger.Info("Starting bookmarks backup.")

	err = r.LoadBookmarks()
	if err != nil {
		return err
	}

	// Get updated and new bookmarks
	newBookmarks, changedBookmarks, removedBookmarks, err := r.GetBookmarkChanges()
	if err != nil {
		return err
	}

	// Download permanent copy where ready and file still missing
	downloadCount := 0
	failedIds := make(map[uint64]bool)

	for _, bookmark := range newBookmarks {
		downloaded, err := r.API.DownloadFileIfMissing(bookmark.Title, bookmark.Id, r.RaindropRoot)
		if err != nil {
			r.Logger.Warn(err.Error())
			failedIds[bookmark.Id] = true
		}

		if downloaded {
			downloadCount++
		}
	}

	for _, bookmark := range changedBookmarks {
		downloaded, err := r.API.DownloadFileIfMissing(bookmark.Title, bookmark.Id, r.RaindropRoot)
		if err != nil {
			r.Logger.Warn(err.Error())
			failedIds[bookmark.Id] = true
		}

		if downloaded {
			downloadCount++
		}
	}

	// Merge unchanged bookmarks with changed/new
	keptIds := make(map[uint64]bool)
	r.UpdatedBookmarks = make([]*data.Bookmark, 0, len(r.Bookmarks)+len(changedBookmarks))

	for _, bookmark := range newBookmarks {
		if _, exists := failedIds[bookmark.Id]; exists {
			continue
		}
		r.UpdatedBookmarks = append(r.UpdatedBookmarks, bookmark)
		keptIds[bookmark.Id] = true
	}

	for _, bookmark := range changedBookmarks {
		if _, exists := failedIds[bookmark.Id]; exists {
			continue
		}
		r.UpdatedBookmarks = append(r.UpdatedBookmarks, bookmark)
		keptIds[bookmark.Id] = true
	}

	for _, bookmark := range r.Bookmarks {
		if _, exists := failedIds[bookmark.Id]; exists {
			continue
		}

		if _, exists := keptIds[bookmark.Id]; exists {
			continue
		}

		if !slices.Contains(removedBookmarks, bookmark.Id) {
			r.UpdatedBookmarks = append(r.UpdatedBookmarks, bookmark)
		}
	}

	// Delete files for removed bookmarks
	for bookmarkId := range removedBookmarks {
		targetDir := filepath.Join(r.RaindropRoot, fmt.Sprintf("%d", bookmarkId))
		r.Logger.Infof("Deleting %s because the bookmark was removed.", targetDir)
		err := os.RemoveAll(targetDir)
		if err != nil {
			r.Logger.Errorf("Failed to remove %s: %s.", targetDir, err.Error())
		}
	}

	// r.PruneBookmarkBackups()
	r.PruneBackupFiles("bookmark")

	if len(newBookmarks) > 0 || len(changedBookmarks) > 0 || len(removedBookmarks) > 0 || downloadCount > 0 {
		err = r.SaveBookmarksBackupFile()
		if err != nil {
			return err
		}
	}

	// Report
	var changedString string = "bookmarks"
	var fileString string = "files"
	var newString string = "bookmarks"
	var removedString = "bookmarks"

	if len(changedBookmarks) == 1 {
		changedString = "bookmark"
	}

	if downloadCount == 1 {
		fileString = "file"
	}

	if len(newBookmarks) == 1 {
		newString = "bookmark"
	}

	if len(removedBookmarks) == 1 {
		removedString = "bookmark"
	}
	r.Logger.Infof(
		"Finished. %d new %s; %d changed %s; %d removed %s; %d %s downloaded.",
		len(newBookmarks),
		newString,
		len(changedBookmarks),
		changedString,
		len(removedBookmarks),
		removedString,
		downloadCount,
		fileString,
	)

	return nil
}

func (r *Raindrop) ListBookmarks() (bookmarks map[uint64]*data.Bookmark, err error) {
	bookmarks, err = r.getAllBookmarks()
	if err != nil {
		return bookmarks, err
	}

	return bookmarks, nil
}

func (r *Raindrop) LoadBookmarks() (err error) {
	bookmarks := make([]*data.Bookmark, 0)

	if util.PathExists(r.BookmarksFile) {
		contents, err := os.ReadFile(r.BookmarksFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(contents, &bookmarks)
		if err != nil {
			return err
		}
	}

	for _, bookmark := range bookmarks {
		r.Bookmarks[bookmark.Id] = bookmark
	}

	return nil
}

func (r *Raindrop) GetBookmarkChanges() (new []*data.Bookmark, changed []*data.Bookmark, removed []uint64, err error) {
	bookmarks, err := r.getAllBookmarks()
	if err != nil {
		return new, changed, removed, err
	}

	// Find new and changed bookmarks
	for _, bookmark := range bookmarks {
		storedBookmark, exists := r.Bookmarks[bookmark.Id]
		if !exists {
			new = append(new, bookmark)
		} else if bookmark.LastUpdate.After(storedBookmark.LastUpdate) {
			changed = append(changed, bookmark)
		}
	}

	// See if any need deleting
	for _, bookmark := range r.Bookmarks {
		_, exists := bookmarks[bookmark.Id]
		if !exists {
			removed = append(removed, bookmark.Id)
		}
	}

	return new, changed, removed, nil
}

func (r *Raindrop) getAllBookmarks() (bookmarks map[uint64]*data.Bookmark, err error) {
	bookmarks = make(map[uint64]*data.Bookmark)
	page := 0

	for {
		listBookmarksResult, err := r.API.ListBookmarks(page)
		if err != nil {
			return bookmarks, err
		}

		for _, bookmark := range listBookmarksResult.Items {
			bookmarks[bookmark.Id] = bookmark
		}

		over := len(listBookmarksResult.Items) < api.PageSize

		if over {
			break
		}
		page += 1
	}

	return bookmarks, nil
}

func (r *Raindrop) SaveBookmarksBackupFile() (err error) {
	yamlBookmarks, err := yaml.Marshal(r.UpdatedBookmarks)
	if err != nil {
		return nil
	}

	if util.PathExists(r.BookmarksFile) {
		backupFilename := filepath.Join(r.RaindropRoot, fmt.Sprintf("bookmarks-%d.yaml", time.Now().Unix()))

		r.Logger.Infof("Copying %s to %s.", r.BookmarksFile, backupFilename)
		r.Logger.Infof("Saving bookmarks to %s.", r.BookmarksFile)

		err = os.Rename(r.BookmarksFile, backupFilename)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %s", r.BookmarksFile, backupFilename, err.Error())
		}
	}

	err = os.WriteFile(r.BookmarksFile, yamlBookmarks, 0600)
	if err != nil {
		return err
	}

	return nil
}

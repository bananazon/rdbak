package raindrop

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/bananazon/rdbak/pkg/data"
	"github.com/bananazon/rdbak/pkg/util"
	"gopkg.in/yaml.v3"
)

const Second = time.Second
const Minute = time.Minute
const Hour = time.Hour
const Day = 24 * time.Hour
const Week = 7 * Day

func (r *Raindrop) Backup(flagPrune bool) (err error) {
	r.PruneOlder = flagPrune

	r.Logger.Info("Starting bookmarks backup.")

	err = r.LoadRaindrops()
	if err != nil {
		return err
	}

	// Get updated and new bookmarks
	newBookmarks, changedBookmarks, removedBookmarks, err := r.GetChanges()
	if err != nil {
		return err
	}

	// Download permanent copy where ready and file still missing
	downloadCount := 0
	failedIds := make(map[uint64]bool)

	for _, bookmark := range newBookmarks {
		downloaded, err := r.API.DownloadFileIfMissing(bookmark.Title, bookmark.Id, r.Config.ExportDir)
		if err != nil {
			r.Logger.Warn(err.Error())
			failedIds[bookmark.Id] = true
		}

		if downloaded {
			downloadCount++
		}
	}

	for _, bookmark := range changedBookmarks {
		downloaded, err := r.API.DownloadFileIfMissing(bookmark.Title, bookmark.Id, r.Config.ExportDir)
		if err != nil {
			r.Logger.Warn(err.Error())
			failedIds[bookmark.Id] = true
		}

		if downloaded {
			downloadCount++
		}
	}

	// Marge unchanged bookmarks with changed/new
	keptIds := make(map[uint64]bool)
	r.UpdatedRaindrops = make([]*data.Raindrop, 0, len(r.Raindrops)+len(changedBookmarks))

	for _, bookmark := range newBookmarks {
		if _, exists := failedIds[bookmark.Id]; exists {
			continue
		}
		r.UpdatedRaindrops = append(r.UpdatedRaindrops, bookmark)
		keptIds[bookmark.Id] = true
	}

	for _, bookmark := range changedBookmarks {
		if _, exists := failedIds[bookmark.Id]; exists {
			continue
		}
		r.UpdatedRaindrops = append(r.UpdatedRaindrops, bookmark)
		keptIds[bookmark.Id] = true
	}

	for _, bookmark := range r.Raindrops {
		if _, exists := failedIds[bookmark.Id]; exists {
			continue
		}

		if _, exists := keptIds[bookmark.Id]; exists {
			continue
		}

		if !slices.Contains(removedBookmarks, bookmark.Id) {
			r.UpdatedRaindrops = append(r.UpdatedRaindrops, bookmark)
		}
	}

	// Delete files for removed bookmarks
	for bookmarkId := range removedBookmarks {
		targetDir := filepath.Join(r.Config.ExportDir, fmt.Sprintf("%d", bookmarkId))
		r.Logger.Infof("Deleting %s because the bookmark was removed.", targetDir)
		err := os.RemoveAll(targetDir)
		if err != nil {
			r.Logger.Errorf("Failed to remove %s: %s.", targetDir, err.Error())
		}
	}

	r.PruneBackups()

	if len(newBookmarks) > 0 || len(changedBookmarks) > 0 || len(removedBookmarks) > 0 || downloadCount > 0 {
		err = r.SaveBookmarks()
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

func (r *Raindrop) PruneBackups() {
	if r.PruneOlder {
		r.Logger.Info("Looking for outdated backup files to prune.")

		pattern := fmt.Sprintf("%s/%s", r.HomePath, "bookmarks-*.yaml")
		timePeriod := 1 * Week

		oldFiles, err := util.FindOldFiles(pattern, timePeriod)
		if err != nil {
			r.Logger.Warn("failed to find outdated backup files")
		} else {
			if len(oldFiles) > 0 {
				var filesString = "files"
				if len(oldFiles) == 1 {
					filesString = "file"
				}
				r.Logger.Infof("I found %d %s that can be pruned.", len(oldFiles), filesString)
				for _, filename := range oldFiles {
					r.Logger.Infof("Pruning %s", filename)
					err = os.Remove(filename)
					if err != nil {
						r.Logger.Warnf("Failed to delete %s", filename)
					}
				}
			} else {
				r.Logger.Info("No outdated backup files found.")
			}
		}
	}
}

func (r *Raindrop) SaveBookmarks() (err error) {
	yamlBookmarks, err := yaml.Marshal(r.UpdatedRaindrops)
	if err != nil {
		return nil
	}

	if util.PathExists(r.Config.BookmarksFile) {
		backupFilename := filepath.Join(r.HomePath, fmt.Sprintf("bookmarks-%d.yaml", time.Now().Unix()))

		r.Logger.Infof("Copying %s to %s.", r.Config.BookmarksFile, backupFilename)
		r.Logger.Infof("Saving bookmarks to %s.", r.Config.BookmarksFile)

		err = os.Rename(r.Config.BookmarksFile, backupFilename)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %s", r.Config.BookmarksFile, backupFilename, err.Error())
		}
	}

	err = os.WriteFile(r.Config.BookmarksFile, yamlBookmarks, 0600)
	if err != nil {
		return err
	}

	return nil
}

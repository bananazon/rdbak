package raindrop

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/gdanko/rdbak/pkg/api"
	"github.com/gdanko/rdbak/pkg/crypt"
	"github.com/gdanko/rdbak/pkg/data"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const Second = time.Second
const Minute = time.Minute
const Hour = time.Hour
const Day = 24 * time.Hour
const Week = 7 * Day

type Config struct {
	Email             string `yaml:"email"`
	Password          string `yaml:"password,omitempty"`
	EncryptedPassword string `yaml:"encryptedPassword,omitempty"`
	BookmarksFile     string `yaml:"bookmarksFile"`
	ExportDir         string `yaml:"exportDir"`
}

type Raindrop struct {
	API          *api.APIClient
	Bookmarks    []*data.Bookmark
	ConfigPath   string
	Config       *Config
	HomePath     string
	Logger       *logrus.Logger
	NewBookmarks []*data.Bookmark
	PruneOlder   bool
}

func New(homePath string, configPath string, pruneOlder bool, logger *logrus.Logger) *Raindrop {
	return &Raindrop{
		API:        api.NewApiClient(logger),
		ConfigPath: configPath,
		Config:     &Config{},
		HomePath:   homePath,
		Logger:     logger,
		PruneOlder: pruneOlder,
	}
}

func (r *Raindrop) ParseConfig() (err error) {
	contents, err := os.ReadFile(r.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read %s", r.ConfigPath)
	}

	if err = yaml.Unmarshal(contents, r.Config); err != nil {
		return fmt.Errorf("failed to parse %s", r.ConfigPath)
	}

	if len(r.Config.Password) == 0 {
		if len(r.Config.EncryptedPassword) > 0 {
			err = r.DecryptPassword()
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("both password and passwordEncrypted fields in the config are empty; please fix this")
		}
	}

	return nil
}

func (r *Raindrop) EncryptPassword() (err error) {
	ciphertext, err := crypt.Encrypt(r.Config.Password)
	if err != nil {
		return err
	}
	r.Config.EncryptedPassword = ciphertext
	return nil
}

func (r *Raindrop) DecryptPassword() (err error) {
	plaintext, err := crypt.Decrypt(r.Config.EncryptedPassword)
	if err != nil {
		return err
	}
	r.Config.Password = plaintext
	return nil
}

func (r *Raindrop) Backup() (err error) {
	r.Logger.Info("Starting bookmarks backup.")

	err = r.LoadBookmarks()
	if err != nil {
		return err
	}

	err = r.API.Login(r.Config.Email, r.Config.Password)
	if err != nil {
		return err
	}

	idToBookmark := make(map[uint64]*data.Bookmark)
	for _, bm := range r.Bookmarks {
		idToBookmark[bm.Id] = bm
	}

	// Get updated and new bookmarks
	changedBookmarks, removedBookmarks, err := r.GetChangedAndRemovedBookmarks(idToBookmark)
	if err != nil {
		return err
	}

	// Download permanent copy where ready and file still missing
	downloadCount := 0
	failedIds := make(map[uint64]bool)
	for _, bookmark := range changedBookmarks {
		if bookmark.Cache.Status != "ready" {
			continue
		}
		downloaded, err := r.API.DownloadFileIfMissing(bookmark.Id, r.Config.ExportDir)
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
	r.NewBookmarks = make([]*data.Bookmark, 0, len(r.Bookmarks)+len(changedBookmarks))
	for _, bm := range changedBookmarks {
		if _, exists := failedIds[bm.Id]; exists {
			continue
		}
		r.NewBookmarks = append(r.NewBookmarks, bm)
		keptIds[bm.Id] = true
	}
	for _, bm := range r.Bookmarks {
		if _, exists := failedIds[bm.Id]; exists {
			continue
		}
		if _, exists := keptIds[bm.Id]; exists {
			continue
		}
		if slices.Contains(removedBookmarks, bm.Id) {
			r.Logger.Info("Skipping!")
		} else {
			r.NewBookmarks = append(r.NewBookmarks, bm)
			keptIds[bm.Id] = true
		}
	}

	r.PruneBackups()

	if len(changedBookmarks) > 0 || len(removedBookmarks) > 0 || downloadCount > 0 {
		err = r.SaveBookmarks()
		if err != nil {
			return err
		}

		var changedString string = "bookmarks"
		var fileString string = "files"
		var removedString = "bookmarks"

		if len(changedBookmarks) == 1 {
			changedString = "bookmark"
		}

		if downloadCount == 1 {
			fileString = "file"
		}

		if len(removedBookmarks) == 1 {
			removedString = "bookmark"
		}
		r.Logger.Infof(
			"Finished. %d %s new or changed; %d %s removed; %d new %s downloaded.",
			len(changedBookmarks),
			changedString,
			len(removedBookmarks),
			removedString,
			downloadCount,
			fileString,
		)
	} else {
		r.Logger.Info("No new or changed bookmarks; no removed bookmarks; no new files downloaded.")
	}

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
	yamlBookmarks, err := yaml.Marshal(r.NewBookmarks)
	if err != nil {
		return nil
	}

	if util.FileExists(r.Config.BookmarksFile) {
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

func (r *Raindrop) LoadBookmarks() (err error) {
	r.Bookmarks = make([]*data.Bookmark, 0)

	if util.FileExists(r.Config.BookmarksFile) {
		contents, err := os.ReadFile(r.Config.BookmarksFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(contents, &r.Bookmarks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Raindrop) GetChangedAndRemovedBookmarks(storedBookmarks map[uint64]*data.Bookmark) (changed []*data.Bookmark, removed []uint64, err error) {
	var raindropBookmarks = make(map[uint64]*data.Bookmark)

	page := 0
	// Load all the bookmarks here first
	for {
		listResult, err := r.API.ListBookmarks(page)
		if err != nil {
			return changed, removed, err
		}

		for _, bookmark := range listResult.Items {
			raindropBookmarks[bookmark.Id] = bookmark
		}

		over := len(listResult.Items) < api.PageSize

		if over {
			break
		}
		page += 1
	}

	// Check for updates
	for _, bookmarkItem := range raindropBookmarks {
		storedBookmark, exists := storedBookmarks[bookmarkItem.Id]
		if !exists {
			changed = append(changed, bookmarkItem)
		} else if bookmarkItem.LastUpdate.After(storedBookmark.LastUpdate) {
			changed = append(changed, bookmarkItem)
		}
	}

	// See if any need deleting
	for _, bookmarkItem := range storedBookmarks {
		_, exists := raindropBookmarks[bookmarkItem.Id]
		if !exists {
			removed = append(removed, bookmarkItem.Id)
		}
	}

	return changed, removed, nil
}

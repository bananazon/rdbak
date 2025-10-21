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
	API              *api.APIClient
	Collections      map[uint64]*data.Collection
	Config           *Config
	ConfigPath       string
	HomePath         string
	Logger           *logrus.Logger
	PruneOlder       bool
	Raindrops        map[uint64]*data.Bookmark
	UpdatedRaindrops []*data.Bookmark
}

func New(homePath string, configPath string, pruneOlder bool, logger *logrus.Logger) (rd *Raindrop, err error) {
	rd = &Raindrop{
		API:        api.NewApiClient(logger),
		ConfigPath: configPath,
		Config:     &Config{},
		HomePath:   homePath,
		Logger:     logger,
		PruneOlder: pruneOlder,
	}
	rd.Raindrops = make(map[uint64]*data.Bookmark)

	err = rd.ParseConfig()
	if err != nil {
		return rd, err
	}

	err = rd.API.Login(rd.Config.Email, rd.Config.Password)
	if err != nil {
		return rd, err
	}

	return rd, nil
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

func (r *Raindrop) Backup(flagPrune bool) (err error) {
	r.PruneOlder = flagPrune

	r.Logger.Info("Starting bookmarks backup.")

	err = r.LoadRaindrops()
	if err != nil {
		return err
	}

	// Get updated and new bookmarks
	newBookmarks, changedBookmarks, removedBookmarks, err := r.GetChangedAndRemovedBookmarks()
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
	r.UpdatedRaindrops = make([]*data.Bookmark, 0, len(r.Raindrops)+len(changedBookmarks))

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

func (r *Raindrop) ListRaindrops() (raindrops map[uint64]*data.Bookmark, err error) {
	raindrops, err = r.getAllBookmarks()
	if err != nil {
		return raindrops, err
	}

	return raindrops, nil
}

func (r *Raindrop) ListCollections() (collections map[uint64]*data.Collection, err error) {
	collections = make(map[uint64]*data.Collection)

	listCollectionsResult, err := r.API.ListCollections()
	if err != nil {
		return collections, err
	}

	for _, collection := range listCollectionsResult.Items {
		collections[collection.Id] = collection
	}

	return collections, nil
}

func (r *Raindrop) LoadRaindrops() (err error) {
	bookmarks := make([]*data.Bookmark, 0)

	if util.PathExists(r.Config.BookmarksFile) {
		contents, err := os.ReadFile(r.Config.BookmarksFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(contents, &bookmarks)
		if err != nil {
			return err
		}
	}

	for _, bookmark := range bookmarks {
		r.Raindrops[bookmark.Id] = bookmark
	}

	return nil
}

func (r *Raindrop) getAllBookmarks() (raindrops map[uint64]*data.Bookmark, err error) {
	raindrops = make(map[uint64]*data.Bookmark)
	page := 0

	for {
		listRaindropsResult, err := r.API.ListRaindrops(page)
		if err != nil {
			return raindrops, err
		}

		for _, bookmark := range listRaindropsResult.Items {
			raindrops[bookmark.Id] = bookmark
		}

		over := len(listRaindropsResult.Items) < api.PageSize

		if over {
			break
		}
		page += 1
	}

	return raindrops, nil
}

func (r *Raindrop) GetChangedAndRemovedBookmarks() (new []*data.Bookmark, changed []*data.Bookmark, removed []uint64, err error) {
	raindrops, err := r.getAllBookmarks()
	if err != nil {
		return new, changed, removed, err
	}

	// Find new and changed bookmarks
	for _, bookmark := range raindrops {
		storedBookmark, exists := r.Raindrops[bookmark.Id]
		if !exists {
			new = append(new, bookmark)
		} else if bookmark.LastUpdate.After(storedBookmark.LastUpdate) {
			changed = append(changed, bookmark)
		}
	}

	// See if any need deleting
	for _, bookmark := range r.Raindrops {
		_, exists := raindrops[bookmark.Id]
		if !exists {
			removed = append(removed, bookmark.Id)
		}
	}

	return new, changed, removed, nil
}

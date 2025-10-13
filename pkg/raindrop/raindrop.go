package raindrop

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdanko/rdbak/pkg/api"
	"github.com/gdanko/rdbak/pkg/crypt"
	"github.com/gdanko/rdbak/pkg/data"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const Day = 24 * time.Hour

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
	Logger       *logrus.Logger
	NewBookmarks []*data.Bookmark
	PruneOlder   bool
}

func New(configPath string, pruneOlder bool, logger *logrus.Logger) *Raindrop {
	return &Raindrop{
		API:        api.NewApiClient(logger),
		ConfigPath: configPath,
		Config:     &Config{},
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
	changedBookmarks, err := r.GetChangedBookmarks(idToBookmark)
	if err != nil {
		return err
	}

	// Download permanent copy where ready and file still missing
	downloadCount := 0
	failedIds := make(map[uint64]bool)
	for _, bm := range changedBookmarks {
		if bm.Cache.Status != "ready" {
			continue
		}
		downloaded, err := r.API.DownloadFileIfMissing(bm.Id, r.Config.ExportDir)
		if err != nil {
			failedIds[bm.Id] = true
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
		r.NewBookmarks = append(r.NewBookmarks, bm)
		keptIds[bm.Id] = true
	}

	r.PruneBookmarks()

	if len(changedBookmarks) > 0 || downloadCount > 0 {
		err = r.SaveBookmarks()
		if err != nil {
			return err
		}

		var bookmarkString string = "bookmarks"
		var fileString string = "files"

		if len(changedBookmarks) == 1 {
			bookmarkString = "bookmark"
		}

		if downloadCount == 1 {
			fileString = "file"
		}
		r.Logger.Infof("Finished. %d %s new or changed; %d new %s downloaded.", len(changedBookmarks), bookmarkString, downloadCount, fileString)
	} else {
		r.Logger.Info("No new or changed bookmarks. No new downloaded files.")
	}

	return nil
}

func (r *Raindrop) PruneBookmarks() {
	if r.PruneOlder {
		r.Logger.Info("Looking for outdated backup files to prune.")

		pattern := fmt.Sprintf("%s/%s", r.Config.ExportDir, "bookmarks-*.yaml")
		timePeriod := 7 * Day

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
		backupFilename := filepath.Join(
			r.Config.ExportDir,
			fmt.Sprintf("bookmarks-%d.yaml", time.Now().Unix()),
		)

		r.Logger.Infof("Copying %s to %s.", r.Config.BookmarksFile, backupFilename)
		r.Logger.Infof("Saving bookmarks to %s.", r.Config.BookmarksFile)

		err = os.Rename(r.Config.BookmarksFile, backupFilename)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %s", r.Config.BookmarksFile, backupFilename, err.Error())
		}
	}

	err = os.WriteFile(r.Config.BookmarksFile, yamlBookmarks, 0644)
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

func (r *Raindrop) GetChangedBookmarks(storedBookmarks map[uint64]*data.Bookmark) (changed []*data.Bookmark, err error) {
	page := 0
	for {
		lr, err := r.API.ListBookmarks(page)
		if err != nil {
			return changed, err
		}

		over := len(lr.Items) < api.PageSize
		for _, itm := range lr.Items {
			if storedBm, exists := storedBookmarks[itm.Id]; !exists {
				changed = append(changed, itm)
			} else if itm.LastUpdate.After(storedBm.LastUpdate) {
				changed = append(changed, itm)
			}
		}
		// DBG
		//over = true
		if over {
			break
		}
		page += 1
	}

	return changed, nil
}

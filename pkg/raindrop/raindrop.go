package raindrop

import (
	"fmt"
	"os"

	"github.com/gdanko/rdbak/pkg/api"
	"github.com/gdanko/rdbak/pkg/crypt"
	"github.com/gdanko/rdbak/pkg/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigPath        string
	Email             string `yaml:"email"`
	Password          string `yaml:"password"`
	EncryptedPassword string `yaml:"encryptedPassword"`
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
}

func New(configPath string, logger *logrus.Logger) *Raindrop {
	return &Raindrop{
		API:        api.NewApiClient(logger),
		ConfigPath: configPath,
		Config:     &Config{},
		Logger:     logger,
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
		return fmt.Errorf("please add your plaintext password to %s", r.ConfigPath)
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

func (r *Raindrop) Backup() (err error) {
	r.Logger.Info("Starting bookmarks backup")
	r.Logger.Infof("Using bookmarks file %s", r.Config.BookmarksFile)

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

	// Save updated bookmarks YAML
	err = r.SaveBookmarks()
	if err != nil {
		return err
	}

	// Report
	var bookmarkString string = "bookmarks"
	var fileString string = "files"

	if len(changedBookmarks) == 1 {
		bookmarkString = "bookmark"
	}

	if downloadCount == 1 {
		fileString = "file"
	}
	r.Logger.Infof("Finished. %d %s new or changed; %d new %s downloaded.", len(changedBookmarks), bookmarkString, downloadCount, fileString)

	return nil
}

func (r *Raindrop) SaveBookmarks() (err error) {
	yamlBookmarks, err := yaml.Marshal(r.NewBookmarks)
	if err != nil {
		return nil
	}

	err = os.WriteFile(r.Config.BookmarksFile, yamlBookmarks, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Raindrop) LoadBookmarks() (err error) {
	r.Bookmarks = make([]*data.Bookmark, 0)

	if _, err = os.Stat(r.Config.BookmarksFile); err == nil {
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

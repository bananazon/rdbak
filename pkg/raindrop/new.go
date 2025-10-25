package raindrop

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bananazon/raindrop/pkg/api"
	"github.com/bananazon/raindrop/pkg/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const Second = time.Second
const Minute = time.Minute
const Hour = time.Hour
const Day = 24 * time.Hour
const Week = 7 * Day

type Config struct {
	EncryptedToken string `yaml:"encryptedApiToken,omitempty"`
	Token          string `yaml:"apiToken,omitempty"`
}

type Raindrop struct {
	API                *api.APIClient
	Bookmarks          map[uint64]*data.Bookmark
	BookmarksFile      string
	Collections        map[uint64]*data.Collection
	CollectionsFile    string
	Config             *Config
	ConfigPath         string
	Logger             *logrus.Logger
	PruneOlder         bool
	RaindropRoot       string
	Token              string
	UpdatedBookmarks   []*data.Bookmark
	UpdatedCollections []*data.Collection
}

func New(raindropRoot string, configPath string, logger *logrus.Logger) (r *Raindrop, err error) {
	r = &Raindrop{
		API:          api.NewApiClient(logger),
		Config:       &Config{},
		ConfigPath:   configPath,
		Logger:       logger,
		PruneOlder:   false,
		RaindropRoot: raindropRoot,
	}

	r.Bookmarks = make(map[uint64]*data.Bookmark)
	r.Collections = make(map[uint64]*data.Collection)

	err = r.ParseConfig()
	if err != nil {
		return r, err
	}

	r.BookmarksFile = filepath.Join(raindropRoot, "bookmarks.yaml")
	r.CollectionsFile = filepath.Join(raindropRoot, "collections.yaml")

	r.API.Token = r.Config.Token

	return r, nil
}

func (r *Raindrop) ParseConfig() (err error) {
	contents, err := os.ReadFile(r.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read %s", r.ConfigPath)
	}

	if err = yaml.Unmarshal(contents, r.Config); err != nil {
		return fmt.Errorf("failed to parse %s", r.ConfigPath)
	}

	if len(r.Config.Token) == 0 {
		if len(r.Config.EncryptedToken) > 0 {
			err = r.DecryptToken()
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("both apiToken and apiTokenEncrypted fields in the config are empty; please fix this")
		}
	}

	return nil
}

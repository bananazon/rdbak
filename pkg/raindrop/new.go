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
	Email             string `yaml:"email"`
	Password          string `yaml:"password,omitempty"`
	EncryptedPassword string `yaml:"encryptedPassword,omitempty"`
}

type Raindrop struct {
	API                *api.APIClient
	Bookmarks          map[uint64]*data.Bookmark
	BookmarksFile      string
	Collections        map[uint64]*data.Collection
	CollectionsFile    string
	Config             *Config
	ConfigPath         string
	RaindropRoot       string
	Logger             *logrus.Logger
	PruneOlder         bool
	UpdatedBookmarks   []*data.Bookmark
	UpdatedCollections []*data.Collection
}

func New(raindropRoot string, configPath string, logger *logrus.Logger) (rd *Raindrop, err error) {
	rd = &Raindrop{
		API:          api.NewApiClient(logger),
		Config:       &Config{},
		ConfigPath:   configPath,
		Logger:       logger,
		PruneOlder:   false,
		RaindropRoot: raindropRoot,
	}
	rd.Bookmarks = make(map[uint64]*data.Bookmark)

	err = rd.ParseConfig()
	if err != nil {
		return rd, err
	}

	rd.BookmarksFile = filepath.Join(raindropRoot, "bookmarks.yaml")
	rd.CollectionsFile = filepath.Join(raindropRoot, "collections.yaml")

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

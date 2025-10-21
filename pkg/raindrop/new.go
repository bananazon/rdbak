package raindrop

import (
	"fmt"
	"os"

	"github.com/gdanko/rdbak/pkg/api"
	"github.com/gdanko/rdbak/pkg/data"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

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
	Raindrops        map[uint64]*data.Raindrop
	UpdatedRaindrops []*data.Raindrop
}

func New(homePath string, configPath string, logger *logrus.Logger) (rd *Raindrop, err error) {
	rd = &Raindrop{
		API:        api.NewApiClient(logger),
		ConfigPath: configPath,
		Config:     &Config{},
		HomePath:   homePath,
		Logger:     logger,
		PruneOlder: false,
	}
	rd.Raindrops = make(map[uint64]*data.Raindrop)

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

package cmd

import (
	"fmt"
	"os"

	"github.com/gdanko/rdbak/pkg/globals"
	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	encryptPasswordCmd = &cobra.Command{
		Use:          "encrypt-password",
		Aliases:      []string{"encypt", "e"},
		Short:        "Replace your plaintext password in the config file with an encrypted string",
		Long:         "Replace your plaintext password in the config file with an encrypted string",
		PreRunE:      encryptPasswordPreRunCmd,
		RunE:         encryptPasswordRunCmd,
		SilenceUsage: false,
	}
)

func init() {
	logLevel = logLevelMap[logLevelStr]
	logger = util.ConfigureLogger(logLevel, flagNoColor)

	err = globals.SetHomeDirectory()
	if err != nil {
		logger.Error(err)
		logger.Exit(2)
	}

	GetFlags(encryptPasswordCmd)
	rootCmd.AddCommand(encryptPasswordCmd)
}

func encryptPasswordPreRunCmd(cmd *cobra.Command, args []string) (err error) {
	rd = *raindrop.New(flagConfigFile, logger)
	err = rd.ParseConfig()
	if err != nil {
		return err
	}

	return nil
}

func encryptPasswordRunCmd(cmd *cobra.Command, args []string) (err error) {
	if len(rd.Config.Password) == 0 {
		return fmt.Errorf("config has no plaintext password")
	}

	if len(rd.Config.EncryptedPassword) > 0 {
		return fmt.Errorf("config already has an encrypted password")
	}

	err = rd.EncryptPassword()
	if err != nil {
		return err
	}

	rd.Config.Password = ""

	configBytes, err := yaml.Marshal(&rd.Config)
	if err != nil {
		return err
	}

	err = os.WriteFile(rd.ConfigPath, configBytes, 0644)
	if err != nil {
		return nil
	}

	return nil
}

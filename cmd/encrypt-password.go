package cmd

import (
	"os"

	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	encryptPasswordCmd = &cobra.Command{
		Use:          "encrypt-password",
		Aliases:      []string{"encrypt", "e"},
		Short:        "Replace your plaintext password in the config file with an encrypted string",
		Long:         "Replace your plaintext password in the config file with an encrypted string",
		PreRun:       encryptPasswordPreRunCmd,
		Run:          encryptPasswordRunCmd,
		SilenceUsage: false,
	}
)

func init() {
	GetEncryptPasswordFlags(encryptPasswordCmd)
	rootCmd.AddCommand(encryptPasswordCmd)
}

func encryptPasswordPreRunCmd(cmd *cobra.Command, args []string) {
	rd = *raindrop.New(flagConfigFile, false, logger)
	err = rd.ParseConfig()
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

func encryptPasswordRunCmd(cmd *cobra.Command, args []string) {
	if len(rd.Config.Password) == 0 {
		logger.Error("config has no plaintext password")
		logger.Exit(1)
	}

	if len(rd.Config.EncryptedPassword) > 0 {
		logger.Info("config already has an encrypted password")
		logger.Exit(0)
	}

	err = rd.EncryptPassword()
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}

	rd.Config.Password = ""

	configBytes, err := yaml.Marshal(&rd.Config)
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}

	err = os.WriteFile(rd.ConfigPath, configBytes, 0644)
	if err != nil {
		logger.Error(err)
		logger.Exit(1)
	}
}

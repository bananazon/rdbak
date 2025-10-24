package cmd

import (
	"os"

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
	RootCmd.AddCommand(encryptPasswordCmd)
}

func encryptPasswordPreRunCmd(cmd *cobra.Command, args []string) {

}

func encryptPasswordRunCmd(cmd *cobra.Command, args []string) {
	if len(RD.Config.Password) == 0 {
		Logger.Error("config has no plaintext password")
		Logger.Exit(1)
	}

	if len(RD.Config.EncryptedPassword) > 0 {
		Logger.Info("config already has an encrypted password")
		Logger.Exit(0)
	}

	err = RD.EncryptPassword()
	if err != nil {
		Logger.Error(err)
		Logger.Exit(1)
	}

	RD.Config.Password = ""

	configBytes, err := yaml.Marshal(&RD.Config)
	if err != nil {
		Logger.Error(err)
		Logger.Exit(1)
	}

	err = os.WriteFile(RD.ConfigPath, configBytes, 0600)
	if err != nil {
		Logger.Error(err)
		Logger.Exit(1)
	}
}

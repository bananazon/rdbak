package cmd

import (
	"github.com/gdanko/rdbak/pkg/globals"
	"github.com/gdanko/rdbak/pkg/raindrop"
	"github.com/gdanko/rdbak/pkg/util"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
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

func encryptPasswordRunCmd(cmd *cobra.Command, args []string) error {
	pretty.Println(rd)
	return nil
}

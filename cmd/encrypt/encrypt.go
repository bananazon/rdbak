package encrypt

import (
	"os"

	"github.com/bananazon/raindrop/pkg/context"
	"github.com/bananazon/raindrop/pkg/raindrop"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewEncryptTokenCmd(ctx *context.AppContext) (c *cobra.Command) {
	c = &cobra.Command{
		Use:     "encrypt-token",
		Aliases: []string{"encrypt", "e"},
		Short:   "Replace your plaintext API token in the config file with an encrypted string",
		PreRun: func(cmdC *cobra.Command, args []string) {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("Failed to initialize raindrop: %s", err.Error())
				ctx.Logger.Exit(1)
			}
			ctx.RD = rd
		},
		Run: func(cmdC *cobra.Command, args []string) {
			if len(ctx.RD.Config.Token) == 0 {
				ctx.Logger.Errorf("config has no plaintext API token")
				ctx.Logger.Exit(1)
			}

			if len(ctx.RD.Config.EncryptedToken) > 0 {
				ctx.Logger.Errorf("config already has an encrypted API token")
				ctx.Logger.Exit(1)
			}

			err := ctx.RD.EncryptToken()
			if err != nil {
				ctx.Logger.Errorf("Failed to encrypt the token: %s", err.Error())
				ctx.Logger.Exit(1)
			}

			ctx.RD.Config.Token = ""

			configBytes, err := yaml.Marshal(&ctx.RD.Config)
			if err != nil {
				ctx.Logger.Errorf("Failed to marshal the config to YAML: %s", err.Error())
				ctx.Logger.Exit(1)
			}

			err = os.WriteFile(ctx.RD.ConfigPath, configBytes, 0600)
			if err != nil {
				ctx.Logger.Errorf("Failed to write the configuration file: %s", err.Error())
				ctx.Logger.Exit(1)
			}

			ctx.Logger.Infof("Successfully wrote the configuration file")
		},
	}

	return c
}

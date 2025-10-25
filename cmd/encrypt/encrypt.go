package encrypt

import (
	"fmt"
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
		PreRunE: func(cmdC *cobra.Command, args []string) error {
			rd, err := raindrop.New(ctx.RaindropHome, ctx.RaindropConfig, ctx.Logger)
			if err != nil {
				ctx.Logger.Println("Failed to initialize raindrop:", err.Error())
				return err
			}
			ctx.RD = rd
			return nil
		},
		RunE: func(cmdC *cobra.Command, args []string) error {
			if len(ctx.RD.Config.Token) == 0 {
				ctx.Logger.Println("config has no plaintext API token")
				return fmt.Errorf("config has no plaintext API token")
			}

			if len(ctx.RD.Config.EncryptedToken) > 0 {
				ctx.Logger.Println("config already has an encrypted API token")
				return fmt.Errorf("config already has an encrypted API token")
			}

			err := ctx.RD.EncryptToken()
			if err != nil {
				ctx.Logger.Println("Failed to encrypt the token:", err.Error())
				return err
			}

			ctx.RD.Config.Token = ""

			configBytes, err := yaml.Marshal(&ctx.RD.Config)
			if err != nil {
				ctx.Logger.Println("Failed to marshal the config to YAML:", err.Error())
				return err
			}

			err = os.WriteFile(ctx.RD.ConfigPath, configBytes, 0600)
			if err != nil {
				ctx.Logger.Println("Failed to write the configuration file:", err.Error())
				return err
			}

			ctx.Logger.Println("Successfully wrote the configuration file")

			return nil
		},
	}

	return c
}

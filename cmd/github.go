package cmd

import (
	"fmt"

	"github.com/alex0ptr/toomani/business"
	"github.com/alex0ptr/toomani/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGithubCmd() *cobra.Command {
	var (
		token string
		owner string
	)

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Use GitHub for mani file generation.",
		Long: `The github command uses github.com to generate a mani file.
It crawls repositories from a provided organization or user and creates a mani.yml file, so that you can use it with mani.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindEnv("token", "GITHUB_TOKEN"); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if vToken := viper.GetString("token"); token == "" && vToken != "" {
				token = vToken
			}

			if token == "" {
				return fmt.Errorf("GitHub token is required: provide via --token flag or GITHUB_TOKEN environment variable")
			}

			fileContent, err := business.NewGenerateRepositoryListing(pkg.NewGitHub(token), writer).WriteManagementFile(business.NewPath(owner), matchPrefix, excludePrefix)
			if err != nil {
				return err
			}
			fmt.Println(fileContent)
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "The GitHub token to use for authentication.")
	_ = viper.BindPFlag("token", cmd.PersistentFlags().Lookup("token"))

	cmd.Flags().StringVarP(&owner, "owner", "", "", "The GitHub organization or user to crawl.")
	_ = cmd.MarkFlagRequired("owner")

	return cmd
}

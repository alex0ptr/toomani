package cmd

import (
	"fmt"

	"github.com/alex0ptr/toomani/business"
	"github.com/alex0ptr/toomani/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGitlabCmd() *cobra.Command {
	var (
		host  string
		token string
		group string
	)
	cmd := &cobra.Command{
		Use:   "gitlab",
		Short: "Use GitLab for mani file generation.",
		Long: `The gitlab uses gitlab.com or a private gitlab instance to generate a mani file.
It crawls repositories from a provided group and creates a mani.yml file, so that you can use it with mani.
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindEnv("host", "GITLAB_HOST"); err != nil {
				return err
			}
			if err := viper.BindEnv("token", "GITLAB_TOKEN"); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if vHost := viper.GetString("host"); host == "gitlab.com" && vHost != "" {
				host = vHost
			}
			if vToken := viper.GetString("token"); token == "" && vToken != "" {
				token = vToken
			}

			if token == "" {
				return fmt.Errorf("GitLab token is required: provide via --token flag or GITLAB_TOKEN environment variable")
			}

			url := fmt.Sprintf("https://%s/api/v4", host)
			fileContent, err := business.NewGenerateRepositoryListing(pkg.NewGitLab(url, token), writer).WriteManagementFile(business.NewPath(group), matchPrefix, excludePrefix)
			if err != nil {
				return err
			}
			fmt.Println(fileContent)
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&host, "host", "h", "gitlab.com", "The hostname of the GitLab instance.")
	cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "The GitLab token to use for authentication.")
	_ = viper.BindPFlag("host", cmd.PersistentFlags().Lookup("host"))
	_ = viper.BindPFlag("token", cmd.PersistentFlags().Lookup("token"))
	cmd.Flags().StringVarP(&group, "group", "g", "", "The GitLab group to crawl.")
	_ = cmd.MarkFlagRequired("group")

	return cmd
}

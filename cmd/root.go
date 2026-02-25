package cmd

import (
	"fmt"
	"slices"

	"github.com/alex0ptr/toomani/business"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	writer        business.ConfigurationWriter
	matchPrefix   []string
	excludePrefix []string
)

func NewRootCmd() *cobra.Command {
	outputOptions := []string{"shell", "mani"}

	cmd := &cobra.Command{
		Use:   "toomani",
		Short: "A tool to bootstrap mani configuration files.",
		Long: `toomani is a CLI tool to create configurations for mani.
It crawls repositories from GitLab or GitHub and creates a mani.yml file, so that you can use it with mani.
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				log.SetLevel(log.DebugLevel)
			}

			output, _ := cmd.Flags().GetString("output")
			if !slices.Contains(outputOptions, output) {
				return fmt.Errorf("output must be one of %v", outputOptions)
			}

			matchPrefix, _ = cmd.Flags().GetStringSlice("match-prefix")
			excludePrefix, _ = cmd.Flags().GetStringSlice("exclude-prefix")

			writer = newWriter(output)
			return nil
		},
	}

	cmd.AddCommand(NewGitlabCmd())
	cmd.AddCommand(NewGithubCmd())
	// Remove the -h help shorthand, as we use it for hostname
	cmd.PersistentFlags().BoolP("help", "", false, "help for this command")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	cmd.PersistentFlags().StringSliceP("match-prefix", "m", []string{}, "Only include repositories with this prefix.")
	cmd.PersistentFlags().StringSliceP("exclude-prefix", "e", []string{}, "Exclude repositories with this prefix.")
	cmd.PersistentFlags().StringP("output", "o", "mani", fmt.Sprintf("The output format. Can be one of: %v", outputOptions))

	return cmd
}

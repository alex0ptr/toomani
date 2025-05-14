package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toomani",
		Short: "A tool to bootstrap mani configuration files.",
		Long: `toomani is a CLI tool to create configurations for mani.
It crawls repositories from GitLab or GitHub and creates a mani.yml file, so that you can use it with mani.
`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				log.SetLevel(log.DebugLevel)
			}
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	cmd.AddCommand(NewGitlabCmd())
	// Remove the -h help shorthand, as we use it for hostname
	cmd.PersistentFlags().BoolP("help", "", false, "help for this command")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")

	return cmd
}

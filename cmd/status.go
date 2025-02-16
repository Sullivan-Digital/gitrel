package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current version and the 5 most recent versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		runStatusCmd(ctx)
		return nil
	},
}

func runStatusCmd(ctx interfaces.GitRelContext) {
	git.ShowStatus(ctx)
}

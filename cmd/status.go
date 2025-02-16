package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current version and the 5 most recent versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}

		git.ShowStatus(ctx, gitCtx)
		return nil
	},
}
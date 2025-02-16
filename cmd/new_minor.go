package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)


var newMinorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Increment the minor version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}

		git.IncrementAndCreateBranch("minor", ctx, gitCtx)
		return nil
	},
}
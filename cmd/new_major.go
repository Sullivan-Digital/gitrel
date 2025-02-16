package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)


var newMajorCmd = &cobra.Command{
	Use:   "major",
	Short: "Increment the major version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}

		git.IncrementAndCreateBranch("major", ctx, gitCtx)
		return nil
	},
}
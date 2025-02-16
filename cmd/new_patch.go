package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)


var newPatchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Increment the patch version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}

		git.IncrementAndCreateBranch("patch", ctx, gitCtx)
		return nil
	},
}
package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)


var newPatchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Increment the patch version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		return runNewPatchCmd(ctx)
	},
}

func runNewPatchCmd(ctx interfaces.GitRelContext) error {
	git.IncrementAndCreateBranch("patch", ctx)
	return nil
}

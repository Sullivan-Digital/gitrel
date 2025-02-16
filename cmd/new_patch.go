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

		runNewPatchCmd(ctx)
		return nil
	},
}

func runNewPatchCmd(ctx interfaces.GitRelContext) {
	git.IncrementAndCreateBranch("patch", ctx)
}

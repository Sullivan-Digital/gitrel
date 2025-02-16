package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)


var newPatchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Increment the patch version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := getCommandContext()
		if err != nil {
			return err
		}

		git.IncrementAndCreateBranch("patch", ctx)
		return nil
	},
}
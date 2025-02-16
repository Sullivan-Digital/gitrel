package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)


var newMajorCmd = &cobra.Command{
	Use:   "major",
	Short: "Increment the major version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := getCommandContext()
		if err != nil {
			return err
		}

		git.IncrementAndCreateBranch("major", ctx)
		return nil
	},
}
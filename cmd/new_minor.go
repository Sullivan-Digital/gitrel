package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)


var newMinorCmd = &cobra.Command{
	Use:   "minor",
	Short: "Increment the minor version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		return runNewMinorCmd(ctx)
	},
}

func runNewMinorCmd(ctx interfaces.GitRelContext) error {
	git.IncrementAndCreateBranch("minor", ctx)
	return nil
}

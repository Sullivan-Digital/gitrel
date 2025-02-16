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

		runNewMinorCmd(ctx)
		return nil
	},
}

func runNewMinorCmd(ctx interfaces.GitRelContext) {
	git.IncrementAndCreateBranch("minor", ctx)
}

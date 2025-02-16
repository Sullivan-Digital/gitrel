package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)


var newMajorCmd = &cobra.Command{
	Use:   "major",
	Short: "Increment the major version of the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		runNewMajorCmd(ctx)
		return nil	
	},
}

func runNewMajorCmd(ctx interfaces.GitRelContext) {
	git.IncrementAndCreateBranch("major", ctx)
}

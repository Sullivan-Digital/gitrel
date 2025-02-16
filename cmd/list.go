package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current release branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		return runListCmd(ctx)
	},
}

func runListCmd(ctx interfaces.GitRelContext) error {
	releaseBranches, err := git.ListReleases(ctx)
	if err != nil {
		return err
	}

	ctx.Output().Println("Current release branches:")
	for _, branch := range releaseBranches {
		if branch.IsLocalOnly() {
			ctx.Output().Println(branch.Version + " (local only)")
		} else {
			ctx.Output().Println(branch.Version)
		}
	}

	return nil
}

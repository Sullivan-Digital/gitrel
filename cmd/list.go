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

		runListCmd(ctx)
		return nil
	},
}

func runListCmd(ctx interfaces.GitRelContext) {
	releaseBranches, err := git.ListReleases(ctx)
	if err != nil {
		ctx.Output().Println(err)
		return
	}

	if len(releaseBranches) == 0 {
		ctx.Output().Println("No release branches found.")
		return
	}

	ctx.Output().Println("Current release branches:")
	for _, branch := range releaseBranches {
		if branch.IsLocalOnly() {
			ctx.Output().Println(branch.Version + " (local only)")
		} else {
			ctx.Output().Println(branch.Version)
		}
	}
}

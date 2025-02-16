package cmd

import (
	"fmt"
	"gitrel/git"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List current release branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}

		releaseBranches, err := git.ListReleases(ctx, gitCtx)
		if err != nil {
			return err
		}

		fmt.Println("Current release branches:")
		for _, branch := range releaseBranches {
			fmt.Println(branch)
		}

		return nil
	},
}

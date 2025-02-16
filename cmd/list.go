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
		ctx, err := getCommandContext()
		if err != nil {
			return err
		}

		releaseBranches, err := git.ListReleases(ctx)
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

package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)

var checkoutLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Checkout the latest release branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}


		git.CheckoutVersion("latest", ctx, gitCtx)
		git.ShowStatus(ctx, gitCtx)
		return nil
	},
}


package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)


var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout a release branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}


		git.CheckoutVersion(args[0], ctx, gitCtx)
		git.ShowStatus(ctx, gitCtx)
		return nil
	},
}

func init() {
	checkoutCmd.AddCommand(checkoutVersionCmd)
	checkoutCmd.AddCommand(checkoutLatestCmd)
}

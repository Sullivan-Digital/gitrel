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
		ctx, err := getCommandContext()
		if err != nil {
			return err
		}

		git.CheckoutVersion(args[0], ctx)
		git.ShowStatus(ctx)
		return nil
	},
}

func init() {
	checkoutCmd.AddCommand(checkoutVersionCmd)
	checkoutCmd.AddCommand(checkoutLatestCmd)
}

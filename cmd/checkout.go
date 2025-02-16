package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)


var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout a release branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}


		return runCheckoutCmd(args, ctx)
	},
}

func init() {
	checkoutCmd.AddCommand(checkoutVersionCmd)
	checkoutCmd.AddCommand(checkoutLatestCmd)
}

func runCheckoutCmd(args []string, ctx interfaces.GitRelContext) error {
	git.CheckoutVersion(args[0], ctx)
	git.ShowStatus(ctx)
	return nil
}

package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout [<version> | latest]",
	Short: "Checkout a release branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		runCheckoutCmd(args, ctx)
		return nil
	},
}

var checkoutVersionCmd = &cobra.Command{
	Use:   "<version>",
	Short: "Checkout the release branch matching the specified version prefix",
}

var checkoutLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Checkout the latest release branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		runCheckoutCmd([]string{"latest"}, ctx)
		return nil
	},
}

func init() {
	checkoutCmd.AddCommand(checkoutVersionCmd)
	checkoutCmd.AddCommand(checkoutLatestCmd)
}

func runCheckoutCmd(args []string, ctx interfaces.GitRelContext) {
	err := git.CheckoutVersion(args[0], ctx)
	if err != nil {
		ctx.Output().Println(err)
	}

	git.ShowStatus(ctx)
}

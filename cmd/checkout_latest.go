package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)

var checkoutLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Checkout the latest release branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		return runCheckoutLatestCmd(ctx)
	},
}

func runCheckoutLatestCmd(ctx interfaces.GitRelContext) error {
	git.CheckoutVersion("latest", ctx)
	git.ShowStatus(ctx)
	return nil
}

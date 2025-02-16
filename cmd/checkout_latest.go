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

		runCheckoutLatestCmd(ctx)
		return nil
	},
}

func runCheckoutLatestCmd(ctx interfaces.GitRelContext) {
	err := git.CheckoutVersion("latest", ctx)
	if err != nil {
		ctx.Output().Println(err)
		return
	}

	git.ShowStatus(ctx)
}

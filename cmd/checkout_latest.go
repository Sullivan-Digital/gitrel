package cmd

import (
	"gitrel/git"

	"github.com/spf13/cobra"
)

var checkoutLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Checkout the latest release branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := getCommandContext()
		if err != nil {
			return err
		}

		git.CheckoutVersion("latest", ctx)
		git.ShowStatus(ctx)
		return nil
	},
}


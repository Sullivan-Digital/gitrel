package cmd

import (
	"github.com/spf13/cobra"
)

var checkoutVersionCmd = &cobra.Command{
	Use:   "<version>",
	Short: "Checkout the release branch matching the specified version prefix",
}

package cmd

import (
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Alias for `gitrel update latest`",
	RunE: func(cmd *cobra.Command, args []string) error {
		args = append([]string{"latest"}, args...)
		return updateCmd.RunE(cmd, args)
	},
}

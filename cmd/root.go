package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	RemoteFlag string
	FetchFlag  bool
	NoFetchFlag  bool
	LocalBranchNameFlag string
	RemoteBranchNameFlag string
)

var rootCmd = &cobra.Command{
	Use:   "gitrel",
	Short: "A tool to manage git release branches",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&RemoteFlag, "remote", "", "Specify the git remote name (overrides config)")
	rootCmd.PersistentFlags().BoolVar(&FetchFlag, "fetch", false, "Fetch from remote before listing branches")
	rootCmd.PersistentFlags().BoolVar(&NoFetchFlag, "no-fetch", false, "Do not fetch from remote before listing branches")
	rootCmd.PersistentFlags().StringVar(&LocalBranchNameFlag, "local-branch-name", "", "Specify the local branch name (overrides config)")
	rootCmd.PersistentFlags().StringVar(&RemoteBranchNameFlag, "remote-branch-name", "", "Specify the remote branch name (overrides config)")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(checkoutCmd)
}
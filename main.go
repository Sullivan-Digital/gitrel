package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName(".gitrelrc")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// Look up the directory tree
	dir, err := os.Getwd()
	if err == nil {
		for {
			viper.AddConfigPath(dir)
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	viper.ReadInConfig()
}

func main() {
	initConfig()

	alwaysFetch := viper.GetBool("alwaysFetch")
	remote := viper.GetString("remote")
	if remote == "" {
		remote = "origin"
	}

	var rootCmd = &cobra.Command{
		Use:   "gitrel",
		Short: "A tool to manage git release branches",
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List current release branches",
		Run: func(cmd *cobra.Command, args []string) {
			fetch, _ := cmd.Flags().GetBool("fetch")
			if alwaysFetch {
				fetch = true
			}
			listReleaseBranches(fetch, remote)
		},
	}
	listCmd.PersistentFlags().Bool("fetch", false, "Fetch from remote before listing branches")

	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new release branch",
	}
	newCmd.PersistentFlags().Bool("fetch", false, "Fetch from remote before creating a new branch")

	var newVersionCmd = &cobra.Command{
		Use:   "<version>",
		Short: "Create a new release branch with the specified version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			version := args[0]
			if !validateSemver(version) {
				fmt.Println("Error: Invalid version format. Please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)")
				os.Exit(1)
			}
			createReleaseBranch(version)
		},
	}

	var newMajorCmd = &cobra.Command{
		Use:   "major",
		Short: "Increment the major version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			fetch, _ := cmd.Flags().GetBool("fetch")
			incrementAndCreateBranch("major", fetch)
		},
	}

	var newMinorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Increment the minor version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			fetch, _ := cmd.Flags().GetBool("fetch")
			incrementAndCreateBranch("minor", fetch)
		},
	}

	var newPatchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Increment the patch version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			fetch, _ := cmd.Flags().GetBool("fetch")
			incrementAndCreateBranch("patch", fetch)
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show the current version and the 5 most recent versions",
		Run: func(cmd *cobra.Command, args []string) {
			fetch, _ := cmd.Flags().GetBool("fetch")
			showStatus(fetch)
		},
	}
	statusCmd.PersistentFlags().Bool("fetch", false, "Fetch from remote before showing status")

	var checkoutCmd = &cobra.Command{
		Use:   "checkout",
		Short: "Checkout a release branch",
	}
	checkoutCmd.PersistentFlags().Bool("fetch", false, "Fetch from remote before checking out")

	var checkoutVersionCmd = &cobra.Command{
		Use:   "<version>",
		Short: "Checkout the release branch matching the specified version prefix",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fetch, _ := cmd.Flags().GetBool("fetch")
			checkoutVersion(args[0], fetch)
		},
	}

	var checkoutLatestCmd = &cobra.Command{
		Use:   "latest",
		Short: "Checkout the latest release branch",
		Run: func(cmd *cobra.Command, args []string) {
			fetch, _ := cmd.Flags().GetBool("fetch")
			checkoutVersion("latest", fetch)
		},
	}

	newCmd.AddCommand(newVersionCmd, newMajorCmd, newMinorCmd, newPatchCmd)
	checkoutCmd.AddCommand(checkoutVersionCmd, checkoutLatestCmd)
	rootCmd.AddCommand(listCmd, newCmd, statusCmd, checkoutCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

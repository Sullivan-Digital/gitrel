package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Function to display help message
func showHelp() {
	fmt.Println("Usage: <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help       Show this help message and exit")
	fmt.Println("  list       List current release branches (fetches from remote first)")
	fmt.Println("  new <ver>  Create a new release branch with the specified version")
	fmt.Println("  major      Increment the major version of the latest release")
	fmt.Println("  minor      Increment the minor version of the latest release")
	fmt.Println("  patch      Increment the patch version of the latest release")
	fmt.Println("  current    Show the highest version from existing release branches")
	fmt.Println("  checkout <ver> Checkout the latest release branch matching the specified version prefix")
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "git-tool",
		Short: "A tool to manage git release branches",
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List current release branches",
		Run: func(cmd *cobra.Command, args []string) {
			listReleaseBranches()
		},
	}

	var newCmd = &cobra.Command{
		Use:   "new <version>",
		Short: "Create a new release branch with the specified version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createReleaseBranch(args[0])
		},
	}

	var majorCmd = &cobra.Command{
		Use:   "major",
		Short: "Increment the major version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			incrementAndCreateBranch("major")
		},
	}

	var minorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Increment the minor version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			incrementAndCreateBranch("minor")
		},
	}

	var patchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Increment the patch version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			incrementAndCreateBranch("patch")
		},
	}

	var currentCmd = &cobra.Command{
		Use:   "current",
		Short: "Show the highest version from existing release branches",
		Run: func(cmd *cobra.Command, args []string) {
			highestVersion := getHighestVersion()
			if highestVersion == "0.0.0" {
				fmt.Println("No existing release branches found.")
			} else {
				fmt.Println("Current highest version:", highestVersion)
			}
		},
	}

	var checkoutCmd = &cobra.Command{
		Use:   "checkout <version>",
		Short: "Checkout the latest release branch matching the specified version prefix",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkoutVersion(args[0])
		},
	}

	rootCmd.AddCommand(listCmd, newCmd, majorCmd, minorCmd, patchCmd, currentCmd, checkoutCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func incrementAndCreateBranch(part string) {
	highestVersion := getHighestVersion()
	newVersion := ""
	if highestVersion == "0.0.0" {
		newVersion = "0.1.0"
	} else {
		newVersion = incrementVersion(highestVersion, part)
	}

	createReleaseBranch(newVersion)
}

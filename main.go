package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gitrel",
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
		Use:   "new",
		Short: "Create a new release branch",
	}

	var newVersionCmd = &cobra.Command{
		Use:   "<version>",
		Short: "Create a new release branch with the specified version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createReleaseBranch(args[0])
		},
	}

	var newMajorCmd = &cobra.Command{
		Use:   "major",
		Short: "Increment the major version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			incrementAndCreateBranch("major")
		},
	}

	var newMinorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Increment the minor version of the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			incrementAndCreateBranch("minor")
		},
	}

	var newPatchCmd = &cobra.Command{
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

	newCmd.AddCommand(newVersionCmd, newMajorCmd, newMinorCmd, newPatchCmd)
	rootCmd.AddCommand(listCmd, newCmd, currentCmd, checkoutCmd)
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

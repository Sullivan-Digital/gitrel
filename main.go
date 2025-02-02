package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

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

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show the current version and the 5 most recent versions",
		Run: func(cmd *cobra.Command, args []string) {
			showStatus()
		},
	}

	var checkoutCmd = &cobra.Command{
		Use:   "checkout",
		Short: "Checkout a release branch",
	}

	var checkoutVersionCmd = &cobra.Command{
		Use:   "<version>",
		Short: "Checkout the release branch matching the specified version prefix",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkoutVersion(args[0])
		},
	}

	var checkoutLatestCmd = &cobra.Command{
		Use:   "latest",
		Short: "Checkout the latest release branch",
		Run: func(cmd *cobra.Command, args []string) {
			checkoutVersion("latest")
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

func showStatus() {
	releaseBranches, err := getRemoteReleaseBranches()
	if err != nil {
		fmt.Println(err)
		return
	}

	var versions []string
	for _, branch := range releaseBranches {
		version := strings.TrimPrefix(branch, "release/")
		if validateSemver(version) {
			versions = append(versions, version)
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return compareSemver(versions[i], versions[j])
	})

	if len(versions) == 0 {
		fmt.Println("No existing release branches found.")
		return
	}

	fmt.Println("Current highest version:", versions[len(versions)-1])
	fmt.Println("Recent versions:")
	for i := len(versions) - 1; i >= 0 && i >= len(versions)-5; i-- {
		fmt.Println(versions[i])
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
)

// Function to display help message
func showHelp() {
	fmt.Println("Usage: [--list | --new <version> | --major | --minor | --patch | --current | --checkout <ver>]")
	fmt.Println()
	fmt.Println("List existing release branches or create and push a new release branch")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help     Show this help message and exit")
	fmt.Println("  --list     List current release branches (fetches from remote first)")
	fmt.Println("  --new      Create a new release branch with the specified version")
	fmt.Println("             Example: --new 1.1.0")
	fmt.Println("  --major    Increment the major version of the latest release")
	fmt.Println("  --minor    Increment the minor version of the latest release")
	fmt.Println("  --patch    Increment the patch version of the latest release")
	fmt.Println("  --current  Show the highest version from existing release branches")
	fmt.Println("  --checkout Checkout the latest release branch matching the specified version prefix")
	fmt.Println("             Example: --checkout 1 or --checkout 2.3")
}

func main() {
	helpFlag := flag.Bool("help", false, "Show help message and exit")
	listFlag := flag.Bool("list", false, "List current release branches")
	newVersion := flag.String("new", "", "Create a new release branch with the specified version")
	majorFlag := flag.Bool("major", false, "Increment the major version of the latest release")
	minorFlag := flag.Bool("minor", false, "Increment the minor version of the latest release")
	patchFlag := flag.Bool("patch", false, "Increment the patch version of the latest release")
	currentFlag := flag.Bool("current", false, "Show the highest version from existing release branches")
	checkoutVersionPrefix := flag.String("checkout", "", "Checkout the latest release branch matching the specified version prefix")

	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	if *listFlag {
		listReleaseBranches()
		return
	}

	if *newVersion != "" {
		createReleaseBranch(*newVersion)
		return
	}

	if *majorFlag || *minorFlag || *patchFlag {
		highestVersion := getHighestVersion()
		if highestVersion == "0.0.0" {
			fmt.Println("Error: No existing release branches found to determine the current version.")
			os.Exit(1)
		}

		var newVersion string
		if *majorFlag {
			newVersion = incrementVersion(highestVersion, "major")
		} else if *minorFlag {
			newVersion = incrementVersion(highestVersion, "minor")
		} else if *patchFlag {
			newVersion = incrementVersion(highestVersion, "patch")
		}

		createReleaseBranch(newVersion)
		return
	}

	if *currentFlag {
		highestVersion := getHighestVersion()
		if highestVersion == "0.0.0" {
			fmt.Println("No existing release branches found.")
		} else {
			fmt.Println("Current highest version:", highestVersion)
		}
		return
	}

	if *checkoutVersionPrefix != "" {
		checkoutVersion(*checkoutVersionPrefix)
		return
	}

	fmt.Println("Error: Invalid or missing argument")
	showHelp()
	os.Exit(1)
}

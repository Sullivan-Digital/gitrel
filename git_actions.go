package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Function to execute a shell command and return its output
func execCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Function to fetch and parse branches
func getReleaseBranches(fetch bool) ([]string, error) {
	if fetch {
		fmt.Println("Fetching from remote...")
		_, err := execCommand("git", "fetch", "--all")
		if err != nil {
			return nil, fmt.Errorf("error fetching from remote: %w", err)
		}
	}

	output, err := execCommand("git", "branch", "-r")
	if err != nil {
		return nil, fmt.Errorf("error listing branches: %w", err)
	}

	branches := strings.Split(output, "\n")
	var releaseBranches []string
	for _, branch := range branches {
		if strings.Contains(branch, "origin/release/") {
			releaseBranches = append(releaseBranches, strings.TrimSpace(strings.Replace(branch, "origin/", "", 1)))
		}
	}

	return releaseBranches, nil
}

// Function to list release branches
func listReleaseBranches(fetch bool) {
	releaseBranches, err := getReleaseBranches(fetch)
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Strings(releaseBranches)
	fmt.Println("Current release branches:")
	for _, branch := range releaseBranches {
		fmt.Println(branch)
	}
}

// Function to check if a branch already exists
func branchExists(version string) bool {
	branchName := "release/" + version
	output, _ := execCommand("git", "branch", "--list", branchName)
	if strings.Contains(output, branchName) {
		return true
	}

	output, _ = execCommand("git", "ls-remote", "--heads", "origin", branchName)
	return strings.Contains(output, branchName)
}

// Function to create a new release branch
func createReleaseBranch(version string) {
	if !validateSemver(version) {
		fmt.Println("Error: Invalid version format. Please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)")
		os.Exit(1)
	}

	if branchExists(version) {
		fmt.Println("Error: Branch release/" + version + " already exists.")
		os.Exit(1)
	}

	fmt.Println("Creating new release branch: release/" + version)
	_, err := execCommand("git", "switch", "-c", "release/"+version)
	if err != nil {
		fmt.Println("Error creating branch:", err)
		return
	}

	_, err = execCommand("git", "push", "origin", "release/"+version)
	if err != nil {
		fmt.Println("Error pushing branch:", err)
		return
	}

	_, err = execCommand("git", "switch", "-")
	if err != nil {
		fmt.Println("Error switching back to previous branch:", err)
	}
}

// Function to get the highest version from release branches
func getHighestVersion(fetch bool) string {
	releaseBranches, err := getReleaseBranches(fetch)
	if err != nil {
		fmt.Println(err)
		return ""
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

	if len(versions) > 0 {
		return versions[len(versions)-1]
	}
	return "0.0.0"
}

// Function to checkout the latest release branch matching the specified version prefix
func checkoutVersion(prefix string, fetch bool) {
	releaseBranches, err := getReleaseBranches(fetch)
	if err != nil {
		fmt.Println(err)
		return
	}

	var matchingVersions []string
	for _, branch := range releaseBranches {
		version := strings.TrimPrefix(branch, "release/")
		if prefix == "latest" || strings.HasPrefix(version, prefix) {
			matchingVersions = append(matchingVersions, version)
		}
	}

	if len(matchingVersions) == 0 {
		fmt.Printf("No release branches found matching prefix: %s\n", prefix)
		return
	}

	sort.Slice(matchingVersions, func(i, j int) bool {
		return compareSemver(matchingVersions[i], matchingVersions[j])
	})

	latestVersion := matchingVersions[len(matchingVersions)-1]
	fmt.Printf("Checking out release branch: release/%s\n", latestVersion)
	_, err = execCommand("git", "checkout", "release/"+latestVersion)
	if err != nil {
		fmt.Println("Error checking out branch:", err)
	}
}

// Function to find the current branch and determine the version
func getCurrentVersionFromBranch() string {
	output, err := execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		fmt.Println("Error finding current branch:", err)
		return "unknown"
	}

	branchName := strings.TrimSpace(output)
	if strings.HasPrefix(branchName, "release/") {
		version := strings.TrimPrefix(branchName, "release/")
		if validateSemver(version) {
			return version
		}
		return "unknown"
	}

	return fmt.Sprintf("(%s)", branchName)
}

// Function to show status
func showStatus(fetch bool) {
	releaseBranches, err := getReleaseBranches(fetch)
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

	currentVersion := getCurrentVersionFromBranch()
	fmt.Println("Current version:", currentVersion)
	fmt.Println("Latest version:", versions[len(versions)-1])
	fmt.Println("Previous versions:")
	const maxVersions = 5
	for i := len(versions) - 2; i >= 0 && i >= len(versions)-maxVersions; i-- {
		fmt.Printf(" - %s\n", versions[i])
	}

	if len(versions) > maxVersions {
		fmt.Printf("(%d more...)\n", len(versions)-maxVersions)
	}
}

// Function to increment and create a new branch
func incrementAndCreateBranch(part string, fetch bool) {
	highestVersion := getHighestVersion(fetch)
	newVersion := ""
	if highestVersion == "0.0.0" {
		newVersion = "0.1.0"
	} else {
		newVersion = incrementVersion(highestVersion, part)
	}

	createReleaseBranch(newVersion)
}

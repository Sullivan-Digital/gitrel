package git

import (
	"fmt"
	"gitrel/context"
	"gitrel/semver"
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

// Function to list release branches
func ListReleases(ctx context.CommandContext) ([]string, error) {
	releaseBranches, err := getReleases(ctx)
	if err != nil {
		return nil, err
	}

	versions := make([]string, len(releaseBranches))	
	for i, release := range releaseBranches {
		versions[i] = release.Version
	}

	return versions, nil
}

// Function to check if a branch already exists
func branchExists(version string, ctx context.CommandContext) bool {
	branchName := replaceInBranchPattern(ctx.GetLocalBranchName(), version)
	output, _ := execCommand("git", "branch", "--list", branchName)
	if strings.Contains(output, branchName) {
		return true
	}

	output, _ = execCommand("git", "ls-remote", "--heads", "origin", branchName)
	return strings.Contains(output, branchName)
}

// Function to create a new release branch
func CreateReleaseBranch(version string, ctx context.CommandContext) error {
	if !semver.ValidateSemver(version) {
		return fmt.Errorf("invalid version format. please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)")
	}

	localBranchName := replaceInBranchPattern(ctx.GetLocalBranchName(), version)
	remoteBranchName := replaceInBranchPattern(ctx.GetRemoteBranchName(), version)

	if branchExists(version, ctx) {
		return fmt.Errorf("branch %s already exists", localBranchName)
	}

	fmt.Printf("Creating new release branch: %s\n", localBranchName)
	_, err := execCommand("git", "switch", "-c", localBranchName)
	if err != nil {
		return err
	}

	_, err = execCommand("git", "push", ctx.GetRemote(), localBranchName+":"+remoteBranchName)
	if err != nil {
		return err
	}

	_, err = execCommand("git", "switch", "-")
	if err != nil {
		return err
	}

	return nil
}

// Function to get the highest version from release branches
func getHighestVersion(ctx context.CommandContext) string {
	releases, err := getReleases(ctx)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var versions []string
	for _, release := range releases {
		version := release.Version
		if semver.ValidateSemver(version) {
			versions = append(versions, version)
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return semver.CompareSemver(versions[i], versions[j])
	})

	if len(versions) > 0 {
		return versions[len(versions)-1]
	}
	return "0.0.0"
}

// Function to checkout the latest release branch matching the specified version prefix
func CheckoutVersion(prefix string, ctx context.CommandContext) {
	releases, err := getReleases(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	var matchingVersions []string
	for _, release := range releases {
		version := release.Version
		if prefix == "latest" || strings.HasPrefix(version, prefix) {
			matchingVersions = append(matchingVersions, version)
		}
	}

	if len(matchingVersions) == 0 {
		fmt.Printf("No release branches found matching prefix: %s\n", prefix)
		return
	}

	sort.Slice(matchingVersions, func(i, j int) bool {
		return semver.CompareSemver(matchingVersions[i], matchingVersions[j])
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
		return ""
	}

	branchName := strings.TrimSpace(output)
	if strings.HasPrefix(branchName, "release/") {
		version := strings.TrimPrefix(branchName, "release/")
		if semver.ValidateSemver(version) {
			return version
		}

		return ""
	}

	return ""
}

// Function to show status
func ShowStatus(ctx context.CommandContext) {
	releases, err := getReleases(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	var versions []string
	for _, release := range releases {
		version := release.Version
		if semver.ValidateSemver(version) {
			versions = append(versions, version)
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return semver.CompareSemver(versions[i], versions[j])
	})

	if len(versions) == 0 {
		fmt.Println("No existing release branches found.")
		fmt.Println("Remote:", ctx.GetRemote())
		return
	}

	currentVersion := getCurrentVersionFromBranch()

	if currentVersion != "" {
		fmt.Println("Current version:", currentVersion)
	}

	fmt.Println("Latest version:", versions[len(versions)-1])
	fmt.Println("Remote:", ctx.GetRemote())
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
func IncrementAndCreateBranch(part string, ctx context.CommandContext) {
	highestVersion := getHighestVersion(ctx)
	newVersion := ""
	if highestVersion == "0.0.0" {
		newVersion = "0.1.0"
	} else {
		newVersion = semver.IncrementVersion(highestVersion, part)
	}

	CreateReleaseBranch(newVersion, ctx)
}

// Function to list git remotes
func GetDefaultRemote() (string, error) {
	output, err := execCommand("git", "remote")
	if err != nil {
		return "", fmt.Errorf("error listing git remotes: %w", err)
	}

	remotes := strings.Split(strings.TrimSpace(output), "\n")

	if len(remotes) == 0 {
		return "", fmt.Errorf("no remotes found")
	}

	if len(remotes) > 1 {
		return "", fmt.Errorf("multiple remotes found, please specify one")
	}

	return remotes[0], nil
}

package git

import (
	"fmt"
	"gitrel/context"
	"gitrel/semver"
	"sort"
	"strings"
)

// Function to list release branches
func ListReleases(ctx context.CommandContext, gitCtx GitContext) ([]string, error) {
	releaseBranches, err := getReleases(ctx, gitCtx)
	if err != nil {
		return nil, err
	}

	versions := make([]string, len(releaseBranches))
	for i, release := range releaseBranches {
		versions[i] = release.Version
	}

	return versions, nil
}

// Function to create a new release branch
func CreateReleaseBranch(version string, ctx context.CommandContext, gitCtx GitContext) error {
	if !semver.ValidateSemver(version) {
		return fmt.Errorf("invalid version format. please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)")
	}

	localBranchName := replaceInBranchPattern(ctx.GetLocalBranchName(), version)
	remoteBranchName := replaceInBranchPattern("remotes/"+ctx.GetRemote()+"/"+ctx.GetRemoteBranchName(), version)

	localExists, err := gitCtx.BranchExists(localBranchName)
	if err != nil {
		return fmt.Errorf("error checking if branch exists: %w", err)
	}

	remoteExists, err := gitCtx.BranchExists(remoteBranchName)
	if err != nil {
		return fmt.Errorf("error checking if branch exists: %w", err)
	}

	if localExists || remoteExists {
		return fmt.Errorf("branch %s already exists", localBranchName)
	}

	fmt.Printf("Creating new release branch: %s\n", localBranchName)
	err = gitCtx.SwitchToNewBranch(localBranchName)
	if err != nil {
		return err
	}

	err = gitCtx.PushBranch(ctx.GetRemote(), localBranchName+":"+remoteBranchName)
	if err != nil {
		return err
	}

	err = gitCtx.SwitchBack()
	if err != nil {
		return err
	}

	return nil
}

// Function to checkout the latest release branch matching the specified version prefix
func CheckoutVersion(prefix string, ctx context.CommandContext, gitCtx GitContext) {
	releases, err := getReleases(ctx, gitCtx)
	if err != nil {
		fmt.Println(err)
		return
	}

	var matchingReleases []*ReleaseInfo
	for _, release := range releases {
		if prefix == "latest" || strings.HasPrefix(release.Version, prefix) {
			matchingReleases = append(matchingReleases, release)
		}
	}

	if len(matchingReleases) == 0 {
		fmt.Printf("No release branches found matching prefix: %s\n", prefix)
		return
	}

	sort.Slice(matchingReleases, func(i, j int) bool {
		return semver.CompareSemver(matchingReleases[i].Version, matchingReleases[j].Version)
	})

	latestRelease := matchingReleases[len(matchingReleases)-1]
	branch := latestRelease.GetFirstLocalBranch()
	if branch != nil {
		fmt.Printf("Checking out release branch: %s\n", branch)
		err = gitCtx.CheckoutBranch(branch.BranchName)
		if err != nil {
			fmt.Println("Error checking out branch:", err)
			return
		}
	}

	remoteBranch := latestRelease.GetFirstRemoteBranch()
	if remoteBranch == nil {
		fmt.Printf("No release branches found for matching version: %s\n", prefix)
		return
	}

	localBranchName := replaceInBranchPattern(ctx.GetLocalBranchName(), latestRelease.Version)
	fmt.Printf("Creating local branch %s from remote branch %s\n", localBranchName, remoteBranch.BranchName)

	err = gitCtx.CheckoutBranch(remoteBranch.BranchName)
	if err != nil {
		fmt.Println("Error checking out branch:", err)
		return
	}

	err = gitCtx.SwitchToNewBranch(localBranchName)
	if err != nil {
		fmt.Println("Error creating local branch:", err)
		return
	}
}

// Function to show status
func ShowStatus(ctx context.CommandContext, gitCtx GitContext) {
	releases, err := getReleases(ctx, gitCtx)
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

	currentVersion := getCurrentVersionFromBranch(ctx, gitCtx)

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
func IncrementAndCreateBranch(part string, ctx context.CommandContext, gitCtx GitContext) {
	highestVersion := getHighestVersion(ctx, gitCtx)
	newVersion := ""
	if highestVersion == "0.0.0" {
		newVersion = "0.1.0"
	} else {
		newVersion = semver.IncrementVersion(highestVersion, part)
	}

	CreateReleaseBranch(newVersion, ctx, gitCtx)
}

// Function to list git remotes
func GetDefaultRemote(gitCtx GitContext) (string, error) {
	remotes, err := gitCtx.ListRemotes()
	if err != nil {
		return "", fmt.Errorf("error listing git remotes: %w", err)
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("no remotes found")
	}

	if len(remotes) > 1 {
		return "", fmt.Errorf("multiple remotes found, please specify one")
	}

	return remotes[0], nil
}

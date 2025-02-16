package git

import (
	"fmt"
	"gitrel/interfaces"
	"gitrel/semver"
	"sort"
	"strings"
)

// Function to list release branches
func ListReleases(ctx interfaces.GitRelContext) ([]*ReleaseInfo, error) {
	return getReleases(ctx)
}

// Function to create a new release branch
func CreateReleaseBranch(version string, ctx interfaces.GitRelContext) error {
	if !semver.ValidateSemver(version) {
		return fmt.Errorf("invalid version format. please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)")
	}

	localBranchName := replaceInBranchPattern(ctx.Options().GetLocalBranchName(), version)
	remoteBranchName := replaceInBranchPattern("remotes/"+ctx.Options().GetRemote()+"/"+ctx.Options().GetRemoteBranchName(), version)

	localExists, err := ctx.Git().BranchExists(localBranchName)
	if err != nil {
		return fmt.Errorf("error checking if branch exists: %w", err)
	}

	remoteExists, err := ctx.Git().BranchExists(remoteBranchName)
	if err != nil {
		return fmt.Errorf("error checking if branch exists: %w", err)
	}

	if localExists || remoteExists {
		return fmt.Errorf("branch %s already exists", localBranchName)
	}

	ctx.Output().Printf("Creating new release branch: %s\n", localBranchName)
	err = ctx.Git().SwitchToNewBranch(localBranchName)
	if err != nil {
		return err
	}

	err = ctx.Git().PushBranch(ctx.Options().GetRemote(), localBranchName+":"+remoteBranchName)
	if err != nil {
		return err
	}

	err = ctx.Git().SwitchBack()
	if err != nil {
		return err
	}

	return nil
}

// Function to checkout the latest release branch matching the specified version prefix
func CheckoutVersion(prefix string, ctx interfaces.GitRelContext) {
	releases, err := getReleases(ctx)
	if err != nil {
		ctx.Output().Println(err)
		return
	}

	var matchingReleases []*ReleaseInfo
	for _, release := range releases {
		if prefix == "latest" || strings.HasPrefix(release.Version, prefix) {
			matchingReleases = append(matchingReleases, release)
		}
	}

	if len(matchingReleases) == 0 {
		ctx.Output().Printf("No release branches found matching prefix: %s\n", prefix)
		return
	}

	sort.Slice(matchingReleases, func(i, j int) bool {
		return semver.CompareSemver(matchingReleases[i].Version, matchingReleases[j].Version)
	})

	latestRelease := matchingReleases[len(matchingReleases)-1]
	branch := latestRelease.GetFirstLocalBranch()
	if branch != nil {
		ctx.Output().Printf("Checking out release branch: %s\n", branch)
		err = ctx.Git().CheckoutBranch(branch.BranchName)
		if err != nil {
			ctx.Output().Println("Error checking out branch:", err)
			return
		}
	}

	remoteBranch := latestRelease.GetFirstRemoteBranch()
	if remoteBranch == nil {
		ctx.Output().Printf("No release branches found for matching version: %s\n", prefix)
		return
	}

	localBranchName := replaceInBranchPattern(ctx.Options().GetLocalBranchName(), latestRelease.Version)
	ctx.Output().Printf("Creating local branch %s from remote branch %s\n", localBranchName, remoteBranch.BranchName)

	err = ctx.Git().CheckoutBranch(remoteBranch.BranchName)
	if err != nil {
		ctx.Output().Println("Error checking out branch:", err)
		return
	}

	err = ctx.Git().SwitchToNewBranch(localBranchName)
	if err != nil {
		ctx.Output().Println("Error creating local branch:", err)
		return
	}
}

// Function to show status
func ShowStatus(ctx interfaces.GitRelContext) {
	releases, err := getReleases(ctx)
	if err != nil {
		ctx.Output().Println(err)
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
		ctx.Output().Println("No existing release branches found.")
		ctx.Output().Println("Remote:", ctx.Options().GetRemote())
		return
	}

	currentVersion := getCurrentVersionFromBranch(ctx)

	if currentVersion != "" {
		ctx.Output().Println("Current version:", currentVersion)
	}

	ctx.Output().Println("Latest version:", versions[len(versions)-1])
	ctx.Output().Println("Remote:", ctx.Options().GetRemote())
	ctx.Output().Println("Previous versions:")
	const maxVersions = 5
	for i := len(versions) - 2; i >= 0 && i >= len(versions)-maxVersions; i-- {
		ctx.Output().Printf(" - %s\n", versions[i])
	}

	if len(versions) > maxVersions {
		ctx.Output().Printf("(%d more...)\n", len(versions)-maxVersions)
	}
}

// Function to increment and create a new branch
func IncrementAndCreateBranch(part string, ctx interfaces.GitRelContext) {
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
func GetDefaultRemote(gitCtx interfaces.GitContext) (string, error) {
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

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

	localBranchName := replaceInBranchPattern(ctx.Command().GetOptLocalBranchName(), version)
	remoteTrackingBranchName := replaceInBranchPattern("remotes/"+ctx.Command().GetOptRemote()+"/"+ctx.Command().GetOptRemoteBranchName(), version)

	localExists, err := ctx.Git().BranchExists(localBranchName)
	if err != nil {
		return fmt.Errorf("error checking if branch exists: %w", err)
	}

	remoteExists, err := ctx.Git().BranchExists(remoteTrackingBranchName)
	if err != nil {
		return fmt.Errorf("error checking if branch exists: %w", err)
	}

	if localExists || remoteExists {
		return fmt.Errorf("branch %s already exists", localBranchName)
	}

	err = ctx.Git().SwitchToNewBranch(localBranchName)
	if err != nil {
		return err
	}

	ctx.Output().Printf("Created new release branch: %s\n", localBranchName)

	remoteBranchName := replaceInBranchPattern(ctx.Command().GetOptRemoteBranchName(), version)

	if localBranchName != remoteBranchName {
		ctx.Output().Printf("Pushing %v to %v (%v)...\n", localBranchName, ctx.Command().GetOptRemote(), remoteBranchName)
	} else {
		ctx.Output().Printf("Pushing %v to %v...\n", localBranchName, ctx.Command().GetOptRemote())
	}

	err = ctx.Git().PushBranch(ctx.Command().GetOptRemote(), localBranchName+":"+remoteBranchName)
	if err != nil {
		return err
	}

	ctx.Output().Println("Pushed!")

	err = ctx.Git().SwitchBack()
	if err != nil {
		return err
	}

	curBranch, err := ctx.Git().GetCurrentBranch()
	if err != nil {
		return err
	}

	ctx.Output().Printf("Switched back to branch: %s\n", curBranch)

	return nil
}

// Function to checkout the latest release branch matching the specified version prefix
func CheckoutVersion(prefix string, ctx interfaces.GitRelContext) error {
	releases, err := getReleases(ctx)
	if err != nil {
		return err
	}

	var matchingReleases []*ReleaseInfo
	for _, release := range releases {
		if prefix == "latest" || strings.HasPrefix(release.Version, prefix) {
			matchingReleases = append(matchingReleases, release)
		}
	}

	if len(matchingReleases) == 0 {
		return fmt.Errorf("no release branches found matching prefix: %s", prefix)
	}

	sort.Slice(matchingReleases, func(i, j int) bool {
		return semver.CompareSemver(matchingReleases[i].Version, matchingReleases[j].Version)
	})

	latestRelease := matchingReleases[len(matchingReleases)-1]
	branch := latestRelease.GetFirstLocalBranch()
	if branch != nil {
		ctx.Output().Printf("Checking out release branch: %s\n", branch.BranchName)
		return ctx.Git().CheckoutBranch(branch.BranchName)
	}

	remoteBranch := latestRelease.GetFirstRemoteBranch()
	if remoteBranch == nil {
		return fmt.Errorf("no release branches found for matching version: %s", prefix)
	}

	localBranchName := replaceInBranchPattern(ctx.Command().GetOptLocalBranchName(), latestRelease.Version)
	ctx.Output().Printf("Creating local branch %s from remote branch %s\n", localBranchName, remoteBranch.BranchName)

	err = ctx.Git().CheckoutBranch(remoteBranch.BranchName)
	if err != nil {
		return err
	}

	err = ctx.Git().SwitchToNewBranch(localBranchName)
	if err != nil {
		return err
	}

	return nil
}

// Function to show status
func ShowStatus(ctx interfaces.GitRelContext) {
	releases, err := getReleases(ctx)
	if err != nil {
		ctx.Output().Println(err)
		return
	}

	if len(releases) == 0 {
		ctx.Output().Println("No existing release branches found.")
		ctx.Output().Println("Remote:", ctx.Command().GetOptRemote())
		return
	}

	currentVersion := getCurrentVersionFromBranch(ctx)
	if currentVersion != "" {
		ctx.Output().Println("Current version:", currentVersion)
	} else {
		ctx.Output().Println("Current version: (not on a release branch)")
	}

	latestVersion := releases[len(releases)-1].Version
	ctx.Output().Println("Latest version:", latestVersion)
	ctx.Output().Println("Remote:", ctx.Command().GetOptRemote())
	ctx.Output().Println("Other versions:")

	type ReleaseMetadata struct {
		Version string
		Tags []string
	}

	releaseMetadata := []ReleaseMetadata{}
	for i := len(releases) - 1; i >= 0; i-- {
		tags := []string{}

		if releases[i].Version == latestVersion {
			tags = append(tags, "latest")
		}
		
		if releases[i].Version == currentVersion {
			tags = append(tags, "current")
		}

		md := ReleaseMetadata{
			Version: releases[i].Version,
			Tags: tags,
		}

		releaseMetadata = append(releaseMetadata, md)
	}

	maxI := len(releaseMetadata) - 1
	skippedVersion := false
	for i, md := range releaseMetadata {
		// Display versions on either side of a tag, and tagged versions themselves
		print := false
		if i < maxI && len(releaseMetadata[i+1].Tags) > 0 {
			print = true
		} else if i > 0 && len(releaseMetadata[i-1].Tags) > 0 {
			print = true
		} else if len(md.Tags) > 0 {
			print = true
		}

		if !print {
			skippedVersion = true
			continue
		}

		if skippedVersion {
			ctx.Output().Println(" - ...")
			skippedVersion = false
		}

		if len(md.Tags) > 0 {
			ctx.Output().Printf(" - %s (%s)\n", md.Version, strings.Join(md.Tags, ", "))
		} else {
			ctx.Output().Printf(" - %s\n", md.Version)
		}
	}

	if skippedVersion {
		ctx.Output().Println(" - ...")
	}
}

// Function to increment and create a new branch
func IncrementAndCreateBranch(part string, ctx interfaces.GitRelContext) {
	highestVersion := getHighestVersion(ctx)
	newVersion := ""
	if highestVersion == "0.0.0" {
		if part == "major" {
			newVersion = "1.0.0"
		} else if part == "patch" {
			newVersion = "0.0.1"
		} else {
			// default to a minor version
			newVersion = "0.1.0"
		}
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

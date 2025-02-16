package git

import (
	"fmt"
	"gitrel/context"
	"gitrel/semver"
	"gitrel/utils"
	"sort"
)

// Function to fetch and parse branches
func getReleases(ctx context.CommandContext, gitCtx GitContext) ([]*ReleaseInfo, error) {
	if ctx.GetFetch() {
		fmt.Printf("Fetching from remote '%s'...\n", ctx.GetRemote())
		err := gitCtx.FetchRemote(ctx.GetRemote())
		if err != nil {
			return nil, fmt.Errorf("error fetching from remote: %w", err)
		}
	}

	branches, err := gitCtx.ListAllBranches()
	if err != nil {
		return nil, fmt.Errorf("error listing branches: %w", err)
	}

	remoteBranchPattern := "remotes/" + ctx.GetRemote() + "/" + ctx.GetRemoteBranchName()
	localBranchPattern := ctx.GetLocalBranchName()

	releaseMap := make(map[string]*ReleaseInfo)
	for _, branch := range branches {
		branchType := ""
		version := ""
		if version = getVersionFromBranch(branch, remoteBranchPattern); version != "" {
			branchType = "remote"
		} else if version = getVersionFromBranch(branch, localBranchPattern); version != "" {
			branchType = "local"
		} else {
			continue
		}

		if _, ok := releaseMap[version]; !ok {
			releaseMap[version] = &ReleaseInfo{
				Version:  version,
				Branches: []ReleaseBranch{},
			}
		}

		info := releaseMap[version]
		info.Branches = append(info.Branches, ReleaseBranch{
			BranchName: branch,
			Type:       branchType,
		})
	}

	releases := utils.MapKeys(releaseMap)
	sort.Slice(releases, func(i, j int) bool {
		return semver.CompareSemver(releases[i], releases[j])
	})

	releaseInfos := make([]*ReleaseInfo, 0, len(releases))
	for _, version := range releases {
		releaseInfos = append(releaseInfos, releaseMap[version])
	}

	return releaseInfos, nil
}

// Function to get the highest version from release branches
func getHighestVersion(ctx context.CommandContext, gitCtx GitContext) string {
	releases, err := getReleases(ctx, gitCtx)
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

// Function to find the current branch and determine the version
func getCurrentVersionFromBranch(ctx context.CommandContext, gitCtx GitContext) string {
	branchName, err := gitCtx.GetCurrentBranch()
	if err != nil {
		fmt.Println("Error finding current branch:", err)
		return ""
	}

	version := getVersionFromBranch(branchName, ctx.GetLocalBranchName())
	if version == "" {
		version = getVersionFromBranch(branchName, ctx.GetRemoteBranchName())
	}

	return version
}

func replaceInBranchPattern(branchPattern string, version string) string {
	return fmt.Sprintf(branchPattern, version)
}

func getVersionFromBranch(branch string, branchPattern string) string {
	var version string
	fmt.Sscanf(branch, branchPattern, &version)
	return version
}

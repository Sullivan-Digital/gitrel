package git

import (
	"fmt"
	"gitrel/interfaces"
	"gitrel/semver"
	"gitrel/utils"
	"sort"
)

// getReleases returns an ordered list of releases
func getReleases(ctx interfaces.GitRelContext) ([]*ReleaseInfo, error) {
	if ctx.Command().GetOptFetch() && !ctx.Command().GetFetched() {
		ctx.Output().Printf("Fetching from remote '%s'...\n", ctx.Command().GetOptRemote())
		err := ctx.Git().FetchRemote(ctx.Command().GetOptRemote())
		if err != nil {
			return nil, fmt.Errorf("error fetching from remote: %w", err)
		}

		ctx.Command().SetFetched(true)
	}

	branches, err := ctx.Git().ListAllBranches()
	if err != nil {
		return nil, fmt.Errorf("error listing branches: %w", err)
	}

	remoteBranchPattern := "remotes/" + ctx.Command().GetOptRemote() + "/" + ctx.Command().GetOptRemoteBranchName()
	localBranchPattern := ctx.Command().GetOptLocalBranchName()

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

		if !semver.ValidateSemver(version) {
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
func getHighestVersion(ctx interfaces.GitRelContext) string {
	releases, err := getReleases(ctx)
	if err != nil {
		ctx.Output().Println(err)
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
func getCurrentVersionFromBranch(ctx interfaces.GitRelContext) string {
	branchName, err := ctx.Git().GetCurrentBranch()
	if err != nil {
		ctx.Output().Println("Error finding current branch:", err)
		return ""
	}

	version := getVersionFromBranch(branchName, ctx.Command().GetOptLocalBranchName())
	if version == "" {
		version = getVersionFromBranch(branchName, ctx.Command().GetOptRemoteBranchName())
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

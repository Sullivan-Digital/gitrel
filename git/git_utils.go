package git

import (
	"fmt"
	"gitrel/context"
	"gitrel/semver"
	"gitrel/utils"
	"sort"
	"strings"
)

// Function to fetch and parse branches
func getReleases(ctx context.CommandContext) ([]*ReleaseInfo, error) {
	if ctx.GetFetch() {
		fmt.Printf("Fetching from remote '%s'...\n", ctx.GetRemote())
		_, err := execCommand("git", "fetch", ctx.GetRemote())
		if err != nil {
			return nil, fmt.Errorf("error fetching from remote: %w", err)
		}
	}

	output, err := execCommand("git", "branch", "-a")
	if err != nil {
		return nil, fmt.Errorf("error listing branches: %w", err)
	}

	remoteBranchPattern := ctx.GetRemote() + "/" + ctx.GetRemoteBranchName()
	localBranchPattern := ctx.GetLocalBranchName()

	branches := strings.Split(output, "\n")
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
			Branch: branch,
			Type:   branchType,
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

func replaceInBranchPattern(branchPattern string, version string) string {
	return fmt.Sprintf(branchPattern, version)
}

func getVersionFromBranch(branch string, branchPattern string) string {
	var version string
	fmt.Sscanf(branch, branchPattern, &version)
	return version
}
package cmd

import (
	"gitrel/gitrel_test"
	"testing"
)

func TestRunStatusCmd_ShowsCurrentAndPreviousVersions(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
		"release/2.0.0",
		"remotes/origin/release/2.0.0",
	}
	ctx.GitContext.CurrentBranch = "release/2.0.0"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly() // No side effects expected
	ctx.OutputContext.AssertOutputLines(
		"Current version: 2.0.0",
		"Latest version: 2.0.0",
		"Remote: origin",
		"Previous versions:",
		" - 1.0.0",
		"(no more versions)",
	)
}

func TestRunStatusCmd_ShowsLatestVersionWhenNotOnReleaseBranch(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
		"release/2.0.0",
		"remotes/origin/release/2.0.0",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly() // No side effects expected
	ctx.OutputContext.AssertOutputLines(
		"Current version: (not on a release branch)",
		"Latest version: 2.0.0",
		"Remote: origin",
		"Previous versions:",
		" - 1.0.0",
		"(no more versions)",
	)
}

func TestRunStatusCmd_PerformsFetchBeforeShowingStatus_IfOptionEnabled(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.Fetch = true
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
	}
	ctx.GitContext.CurrentBranch = "release/1.0.0"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectFetchRemote("origin"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Fetching from remote 'origin'...",
		"Current version: 1.0.0",
		"Latest version: 1.0.0",
		"Remote: origin",
		"Previous versions:",
		"(no more versions)",
	)
}

func TestRunStatusCmd_ShowsCorrectVersionsWithDifferentBranchNamingConvention(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.LocalBranchName = "v/%v"
	ctx.GitContext.Branches = []string{
		"main",
		"v/1.0.0",
		"v/2.0.0",
	}
	ctx.GitContext.CurrentBranch = "v/2.0.0"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly() // No side effects expected
	ctx.OutputContext.AssertOutputLines(
		"Current version: 2.0.0",
		"Latest version: 2.0.0",
		"Remote: origin",
		"Previous versions:",
		" - 1.0.0",
		"(no more versions)",
	)
}

func TestRunStatusCmd_ShowsErrorWhenNoReleasesExist(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly() // No side effects expected
	ctx.OutputContext.AssertOutputLines(
		"No existing release branches found.",
		"Remote: origin",
	)
}

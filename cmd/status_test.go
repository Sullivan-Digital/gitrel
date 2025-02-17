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
		"Other versions:",
		" - 2.0.0 (latest, current)",
		" - 1.0.0",
	)
}

func TestRunStatusCmd_PrintsVersionsEitherSideOfCurrent(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"release/2.0.0",
		"release/3.0.0",
		"release/4.0.0",
		"release/5.0.0",
	}
	ctx.GitContext.CurrentBranch = "release/3.0.0"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly() // No side effects expected
	ctx.OutputContext.AssertOutputLines(
		"Current version: 3.0.0",
		"Latest version: 5.0.0",
		"Remote: origin",
		"Other versions:",
		" - 5.0.0 (latest)",
		" - 4.0.0",
		" - 3.0.0 (current)",
		" - 2.0.0",
	)
}

func TestRunStatusCmd_ShowsSkippedVersions(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"release/1.0.0",
		"release/2.0.0",
		"release/3.0.0",
		"release/4.0.0",
		"release/5.0.0",
		"release/6.0.0",
		"release/7.0.0",
	}
	ctx.GitContext.CurrentBranch = "release/2.0.0"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly() // No side effects expected
	ctx.OutputContext.AssertOutputLines(
		"Current version: 2.0.0",
		"Latest version: 7.0.0",
		"Remote: origin",
		"Other versions:",
		" - 7.0.0 (latest)",
		" - 6.0.0",
		" - ...",
		" - 3.0.0",
		" - 2.0.0 (current)",
		" - 1.0.0",
	)
}

func TestRunStatusCmd_ShowsDoubleSkippedVersions(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"release/1.0.0",
		"release/2.0.0",
		"release/3.0.0",
		"release/4.0.0",
		"release/5.0.0",
		"release/6.0.0",
		"release/7.0.0",
		"release/8.0.0",
		"release/9.0.0",
		"release/10.0.0",
	}
	ctx.GitContext.CurrentBranch = "release/5.0.0"

	// Act
	runStatusCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly() // No side effects expected
	ctx.OutputContext.AssertOutputLines(
		"Current version: 5.0.0",
		"Latest version: 10.0.0",
		"Remote: origin",
		"Other versions:",
		" - 10.0.0 (latest)",
		" - 9.0.0",
		" - ...",
		" - 6.0.0",
		" - 5.0.0 (current)",
		" - 4.0.0",
		" - ...",
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
		"Other versions:",
		" - 2.0.0 (latest)",
		" - 1.0.0",
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
		"Other versions:",
		" - 1.0.0 (latest, current)",
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
		"Other versions:",
		" - 2.0.0 (latest, current)",
		" - 1.0.0",
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

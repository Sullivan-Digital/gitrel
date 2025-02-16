package cmd

import (
	"gitrel/gitrel_test"
	"testing"
)

func TestRunListCmd_PrintsAllReleaseBranches(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
		"release/2.0.0",
		"release/3.0.0",
		"remotes/origin/release/3.0.0",
	}

	// Act
	runListCmd(ctx)

	// Assert
	ctx.OutputContext.AssertOutputLines(
		"Current release branches:",
		"1.0.0",
		"2.0.0 (local only)",
		"3.0.0",
	)
}

func TestRunListCmd_PrintsEmptyWhenNoReleaseBranches(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}

	// Act
	runListCmd(ctx)

	// Assert
	ctx.OutputContext.AssertOutputLines(
		"No release branches found.",
	)
}

func TestRunListCmd_HandlesDifferentBranchNamingConvention(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.LocalBranchName = "v/%v"
	ctx.CommandContext.RemoteBranchName = "v/%v"
	ctx.GitContext.Branches = []string{
		"main",
		"v/1.0.0",
		"remotes/origin/v/2.0.0",
	}

	// Act
	runListCmd(ctx)

	// Assert
	ctx.OutputContext.AssertOutputLines(
		"Current release branches:",
		"1.0.0 (local only)",
		"2.0.0",
	)
}

func TestRunListCmd_HandlesDifferentBranchNamingConvention_BetweenRemoteAndLocal(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.LocalBranchName = "local/v/%v"
	ctx.CommandContext.RemoteBranchName = "v/%v"
	ctx.GitContext.Branches = []string{
		"main",
		"local/v/1.0.0",
		"remotes/origin/v/2.0.0",
	}

	// Act
	runListCmd(ctx)

	// Assert
	ctx.OutputContext.AssertOutputLines(
		"Current release branches:",
		"1.0.0 (local only)",
		"2.0.0",
	)
}

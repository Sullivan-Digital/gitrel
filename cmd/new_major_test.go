package cmd

import (
	"gitrel/gitrel_test"
	"testing"
)

func TestRunNewMajorCmd_IncrementsMajorVersion(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
	}

	// Act
	runNewMajorCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("release/2.0.0"),
		gitrel_test.EffectCheckoutBranch("release/2.0.0"),
		gitrel_test.EffectPushBranch("origin", "release/2.0.0"),
		gitrel_test.EffectSwitchBack(),
	)
	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: release/2.0.0",
		"Pushing release/2.0.0 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunNewMajorCmd_IncrementsMajorVersion_FromNoPreviousReleases(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}

	// Act
	runNewMajorCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("release/1.0.0"),
		gitrel_test.EffectCheckoutBranch("release/1.0.0"),
		gitrel_test.EffectPushBranch("origin", "release/1.0.0"),
		gitrel_test.EffectSwitchBack(),
	)

	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: release/1.0.0",
		"Pushing release/1.0.0 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunNewMajorCmd_IncrementsMajorVersion_WithDifferentLocalBranchNamingConvention(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.LocalBranchName = "v/%v"
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"v/1.0.0",
		"remotes/origin/v/1.0.0",
	}

	// Act
	runNewMajorCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("v/2.0.0"),
		gitrel_test.EffectCheckoutBranch("v/2.0.0"),
		gitrel_test.EffectPushBranch("origin", "v/2.0.0:release/2.0.0"),
		gitrel_test.EffectSwitchBack(),
	)
	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: v/2.0.0",
		"Pushing v/2.0.0 to origin (release/2.0.0)...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

package cmd

import (
	"gitrel/gitrel_test"
	"testing"
)

func TestRunNewPatchCmd_IncrementsPatchVersion(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runNewPatchCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("release/1.0.1"),
		gitrel_test.EffectCheckoutBranch("release/1.0.1"),
		gitrel_test.EffectPushBranch("origin", "release/1.0.1"),
		gitrel_test.EffectCheckoutBranch("main"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: release/1.0.1",
		"Pushing release/1.0.1 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunNewPatchCmd_IncrementsPatchVersion_FromNoPreviousReleases(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runNewPatchCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("release/0.0.1"),
		gitrel_test.EffectCheckoutBranch("release/0.0.1"),
		gitrel_test.EffectPushBranch("origin", "release/0.0.1"),
		gitrel_test.EffectCheckoutBranch("main"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: release/0.0.1",
		"Pushing release/0.0.1 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunNewPatchCmd_IncrementsPatchVersion_WithDifferentLocalBranchNamingConvention(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.LocalBranchName = "v/%v"
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"v/1.0.0",
		"remotes/origin/v/1.0.0",
	}
	ctx.GitContext.CurrentBranch = "main"
	
	// Act
	runNewPatchCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("v/1.0.1"),
		gitrel_test.EffectCheckoutBranch("v/1.0.1"),
		gitrel_test.EffectPushBranch("origin", "v/1.0.1:release/1.0.1"),
		gitrel_test.EffectCheckoutBranch("main"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: v/1.0.1",
		"Pushing v/1.0.1 to origin (release/1.0.1)...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

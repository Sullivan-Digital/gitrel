package cmd

import (
	"gitrel/gitrel_test"
	"testing"
)

func TestRunUpdateCmd_PushesToExistingReleaseBranch(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"release/1.0.0",
		"release/2.0.0",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runUpdateCmd([]string{"2.0.0"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCheckoutBranch("release/2.0.0"),
		gitrel_test.EffectMergeBranch("main"),
		gitrel_test.EffectPushBranch("origin", "release/2.0.0"),
		gitrel_test.EffectCheckoutBranch("main"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Checking out release/2.0.0...",
		"Merging main into release/2.0.0...",
		"Pushing release/2.0.0 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunUpdateCmd_PushesToLatestReleaseBranch_WhenLatestVersionIsSpecified(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"release/1.0.0",
		"release/2.0.0",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runUpdateCmd([]string{"latest"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCheckoutBranch("release/2.0.0"),
		gitrel_test.EffectMergeBranch("main"),
		gitrel_test.EffectPushBranch("origin", "release/2.0.0"),
		gitrel_test.EffectCheckoutBranch("main"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Checking out release/2.0.0...",
		"Merging main into release/2.0.0...",
		"Pushing release/2.0.0 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunUpdateCmd_PrintsErrorWhenUncommittedChanges(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.HasUncommittedChangesFl = true

	// Act
	runUpdateCmd([]string{"2.0.0"}, ctx)

	// Assert
	ctx.GitContext.AssertNoSideEffects()
	ctx.OutputContext.AssertOutputLines(
		"you have uncommitted changes. please commit or stash them before updating",
	)
}

func TestRunUpdateCmd_WithPrereleaseVersion(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"release/1.0.0-beta.1",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runUpdateCmd([]string{"1.0.0-beta.1"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCheckoutBranch("release/1.0.0-beta.1"),
		gitrel_test.EffectMergeBranch("main"),
		gitrel_test.EffectPushBranch("origin", "release/1.0.0-beta.1"),
		gitrel_test.EffectCheckoutBranch("main"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Checking out release/1.0.0-beta.1...",
		"Merging main into release/1.0.0-beta.1...",
		"Pushing release/1.0.0-beta.1 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunUpdateCmd_WithBuildMetadata(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"release/1.0.0+build.1",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runUpdateCmd([]string{"1.0.0+build.1"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCheckoutBranch("release/1.0.0+build.1"),
		gitrel_test.EffectMergeBranch("main"),
		gitrel_test.EffectPushBranch("origin", "release/1.0.0+build.1"),
		gitrel_test.EffectCheckoutBranch("main"),
	)
	ctx.OutputContext.AssertOutputLines(
		"Checking out release/1.0.0+build.1...",
		"Merging main into release/1.0.0+build.1...",
		"Pushing release/1.0.0+build.1 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunUpdateCmd_WithInvalidVersion(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}
	ctx.GitContext.CurrentBranch = "main"

	// Act
	runUpdateCmd([]string{"invalid-version"}, ctx)

	// Assert
	ctx.GitContext.AssertNoSideEffects()
	ctx.OutputContext.AssertOutputLines(
		"invalid version format. please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)",
	)
}

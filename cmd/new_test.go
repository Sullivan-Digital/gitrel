package cmd

import (
	"gitrel/gitrel_test"
	"testing"
)

func TestRunNewCmd_CreatesNewReleaseBranch(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}

	// Act
	runNewCmd([]string{"1.0.0"}, ctx)

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

func TestRunNewCmd_WithPrereleaseVersion(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}

	// Act
	runNewCmd([]string{"1.0.0-beta.1"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("release/1.0.0-beta.1"),
		gitrel_test.EffectCheckoutBranch("release/1.0.0-beta.1"),
		gitrel_test.EffectPushBranch("origin", "release/1.0.0-beta.1"),
		gitrel_test.EffectSwitchBack(),
	)
	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: release/1.0.0-beta.1",
		"Pushing release/1.0.0-beta.1 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunNewCmd_WithBuildMetadata(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}

	// Act
	runNewCmd([]string{"1.0.0+build.1"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectCreateBranch("release/1.0.0+build.1"),
		gitrel_test.EffectCheckoutBranch("release/1.0.0+build.1"),
		gitrel_test.EffectPushBranch("origin", "release/1.0.0+build.1"),
		gitrel_test.EffectSwitchBack(),
	)
	ctx.OutputContext.AssertOutputLines(
		"Created new release branch: release/1.0.0+build.1",
		"Pushing release/1.0.0+build.1 to origin...",
		"Pushed!",
		"Switched back to branch: main",
	)
}

func TestRunNewCmd_WithInvalidVersion(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}

	// Act
	runNewCmd([]string{"invalid-version"}, ctx)

	// Assert
	ctx.GitContext.AssertNoSideEffects()
	ctx.OutputContext.AssertOutputLines(
		"invalid version format. please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)",
	)
}
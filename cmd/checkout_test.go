package cmd

import (
	"gitrel/git"
	"gitrel/gitrel_test"
	"testing"
)

func TestRunCheckoutCmd_ChecksOutSpecifiedVersion_Exact(t *testing.T) {
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

	// Act
	runCheckoutCmd([]string{"2.0.0"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(gitrel_test.EffectCheckoutBranch("release/2.0.0"))
	ctx.OutputContext.AssertOutputLines(
		"Checking out release branch: release/2.0.0",
		gitrel_test.GetStdOutIgnoreSideEffects(ctx, func(ctx2 *gitrel_test.TestGitRelContext) {
			git.ShowStatus(ctx2)
		}),
	)
}

func TestRunCheckoutCmd_ChecksOutLatestVersion_Minor(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.3.10",
		"release/1.10.0",
		"release/1.3.9",
		"release/1.2.500",
	}

	// Act
	runCheckoutCmd([]string{"1.3"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(gitrel_test.EffectCheckoutBranch("release/1.3.10"))
	ctx.OutputContext.AssertOutputLines(
		"Checking out release branch: release/1.3.10",
		gitrel_test.GetStdOutIgnoreSideEffects(ctx, func(ctx2 *gitrel_test.TestGitRelContext) {
			git.ShowStatus(ctx2)
		}),
	)
}

func TestRunCheckoutCmd_ChecksOutLatestVersion_Major(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/2.3.10",
		"release/2.10.0",
		"release/2.3.9",
		"release/2.2.500",
		"release/10.0.1",
	}

	// Act
	runCheckoutCmd([]string{"10"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(gitrel_test.EffectCheckoutBranch("release/10.0.1"))
	ctx.OutputContext.AssertOutputLines(
		"Checking out release branch: release/10.0.1",
		gitrel_test.GetStdOutIgnoreSideEffects(ctx, func(ctx2 *gitrel_test.TestGitRelContext) {
			git.ShowStatus(ctx2)
		}),
	)
}

func TestRunCheckoutCmd_PerformsFetchBeforeCheckingOut_IfOptionEnabled(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.Fetch = true
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
		"release/2.0.0",
		"remotes/origin/release/2.0.0",
	}

	// Act
	runCheckoutCmd([]string{"2.0.0"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(
		gitrel_test.EffectFetchRemote("origin"),
		gitrel_test.EffectCheckoutBranch("release/2.0.0"),
	)
}

func TestRunCheckoutCmd_PrintsErrorIfVersionNotFound(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
	}

	// Act
	runCheckoutCmd([]string{"2.0.0"}, ctx)

	// Assert
	if len(ctx.GitContext.SideEffects) != 0 {
		t.Fatalf("expected 0 side effects, got %v", len(ctx.GitContext.SideEffects))
	}

	ctx.OutputContext.AssertOutputLines(
		"no release branches found matching prefix: 2.0.0",
		gitrel_test.GetStdOutIgnoreSideEffects(ctx, func(ctx2 *gitrel_test.TestGitRelContext) {
			git.ShowStatus(ctx2)
		}),
	)
}

func TestRunCheckoutCmd_ChecksOutSpecifiedVersion_WithDifferentBranchNamingConvention(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.CommandContext.LocalBranchName = "v/%v"
	ctx.GitContext.Branches = []string{
		"main",
		"v/1.0.0",
		"v/2.0.0",
	}

	// Act
	runCheckoutCmd([]string{"2.0.0"}, ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(gitrel_test.EffectCheckoutBranch("v/2.0.0"))
	ctx.OutputContext.AssertOutputLines(
		"Checking out release branch: v/2.0.0",
		gitrel_test.GetStdOutIgnoreSideEffects(ctx, func(ctx2 *gitrel_test.TestGitRelContext) {
			git.ShowStatus(ctx2)
		}),
	)
}

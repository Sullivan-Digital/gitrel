package cmd

import (
	"gitrel/gitrel_test"
	"testing"
)

func TestCheckoutLatest_PrintsErrorIfNoReleases(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
	}

	// Act
	runCheckoutLatestCmd(ctx)

	// Assert
	if len(ctx.GitContext.SideEffects) != 0 {
		t.Fatalf("expected 0 side effects, got %v", len(ctx.GitContext.SideEffects))
	}

	expectedOutput := "no release branches found matching prefix: latest\n"
	ctx.OutputContext.AssertOutput(expectedOutput)
}

func TestCheckoutLatest_ChecksOutLatestRelease(t *testing.T) {
	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext(t)
	ctx.GitContext.Branches = []string{
		"main",
		"remotes/origin/main",
		"release/1.0.0",
		"remotes/origin/release/1.0.0",
		"release/3.0.0",
		"remotes/origin/release/3.0.0",
		"release/2.0.0",
		"remotes/origin/release/2.0.0",
	}

	// Act
	runCheckoutLatestCmd(ctx)

	// Assert
	ctx.GitContext.AssertSideEffectsAreExactly(gitrel_test.EffectCheckoutBranch("release/3.0.0"))
	ctx.OutputContext.AssertOutputLines(
		"Checking out release branch: release/3.0.0",
        "Current version: 3.0.0",
        "Latest version: 3.0.0",
        "Remote: origin",
        "Previous versions:",
        " - 2.0.0",
        " - 1.0.0",
		"(no more versions)",
	)
}

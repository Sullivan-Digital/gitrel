package git

import (
	"gitrel/gitrel_test"
	"slices"
	"testing"
)

func TestListReleases_ReturnsListOfReleases(t *testing.T) {

	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext()
	ctx.CommandContext.Fetch = false

	// Act
	releases, err := ListReleases(ctx)
	if err != nil {
		t.Fatalf("error listing releases: %v", err)
	}

	// Assert
	if len(ctx.GitContext.SideEffects) != 0 {
		t.Fatalf("expected 0 side effects, got %v", len(ctx.GitContext.SideEffects))
	}

	expectedReleases := []string{
		"1.0.0",
		"1.0.1",
		"1.0.2",
		"1.1.0",
		"1.1.1",
		"2.0.0",
	}

	actualReleases := make([]string, len(releases))
	for i, release := range releases {
		actualReleases[i] = release.Version

		if release.IsLocalOnly() {
			actualReleases[i] = release.Version + " (local only)"
		}
	}

	if !slices.Equal(actualReleases, expectedReleases) {
		t.Fatalf("expected %v, got %v", expectedReleases, actualReleases)
	}
}

func TestListReleases_IncludesLocalOnlyReleases(t *testing.T) {

	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext()
	ctx.CommandContext.Fetch = false

	ctx.GitContext.Branches = append(ctx.GitContext.Branches, "release/3.0.0")

	// Act
	releases, err := ListReleases(ctx)
	if err != nil {
		t.Fatalf("error listing releases: %v", err)
	}

	// Assert
	if len(ctx.GitContext.SideEffects) != 0 {
		t.Fatalf("expected 0 side effects, got %v", len(ctx.GitContext.SideEffects))
	}

	expectedReleases := []string{
		"1.0.0",
		"1.0.1",
		"1.0.2",
		"1.1.0",
		"1.1.1",
		"2.0.0",
		"3.0.0 (local only)",
	}

	actualReleases := make([]string, len(releases))
	for i, release := range releases {
		actualReleases[i] = release.Version

		if release.IsLocalOnly() {
			actualReleases[i] = release.Version + " (local only)"
		}
	}

	if !slices.Equal(actualReleases, expectedReleases) {
		t.Fatalf("expected %v, got %v", expectedReleases, actualReleases)
	}
}

func TestListReleases_ShowRemoteOnlyReleasesButDontMark(t *testing.T) {

	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext()
	ctx.CommandContext.Fetch = false

	ctx.GitContext.Branches = append(ctx.GitContext.Branches, "remotes/origin/release/3.0.0")

	// Act
	releases, err := ListReleases(ctx)
	if err != nil {
		t.Fatalf("error listing releases: %v", err)
	}

	// Assert
	if len(ctx.GitContext.SideEffects) != 0 {
		t.Fatalf("expected 0 side effects, got %v", len(ctx.GitContext.SideEffects))
	}

	expectedReleases := []string{
		"1.0.0",
		"1.0.1",
		"1.0.2",
		"1.1.0",
		"1.1.1",
		"2.0.0",
		"3.0.0",
	}

	actualReleases := make([]string, len(releases))
	for i, release := range releases {
		actualReleases[i] = release.Version

		if release.IsLocalOnly() {
			actualReleases[i] = release.Version + " (local only)"
		}
	}

	if !slices.Equal(actualReleases, expectedReleases) {
		t.Fatalf("expected %v, got %v", expectedReleases, actualReleases)
	}
}

func TestListReleases_FetchFromRemoteBeforeListing(t *testing.T) {

	// Arrange
	ctx := gitrel_test.DefaultTestGitRelContext()
	ctx.CommandContext.Fetch = true

	// Act
	_, err := ListReleases(ctx)
	if err != nil {
		t.Fatalf("error listing releases: %v", err)
	}

	// Assert
	if len(ctx.GitContext.SideEffects) != 1 {
		t.Fatalf("expected 1 side effect, got %v", len(ctx.GitContext.SideEffects))
	}

	if ctx.GitContext.SideEffects[0] != gitrel_test.EffectFetchRemote("origin") {
		t.Fatalf("expected fetch remote action")
	}
}

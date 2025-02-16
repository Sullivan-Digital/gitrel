package gitrel_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

type TestGitContext struct {
	Branches       []string
	CurrentBranch  string
	PreviousBranch string
	Remotes        []string
	SideEffects    []TestGitSideEffect
	testCtx        *testing.T
}

func DefaultTestGitContext(t *testing.T) *TestGitContext {
	return &TestGitContext{
		Branches: []string{
			"main",
			"release/1.0.0",
			"release/1.0.1",
			"release/1.0.2",
			"release/1.1.0",
			"release/1.1.1",
			"release/2.0.0",
			"remotes/origin/main",
			"remotes/origin/release/1.0.0",
			"remotes/origin/release/1.0.1",
			"remotes/origin/release/1.0.2",
			"remotes/origin/release/1.1.0",
			"remotes/origin/release/1.1.1",
			"remotes/origin/release/2.0.0",
		},
		CurrentBranch:  "main",
		PreviousBranch: "",
		Remotes:        []string{"origin"},
		SideEffects:    []TestGitSideEffect{},
		testCtx:        t,
	}
}

func (c *TestGitContext) FetchRemote(remote string) error {
	c.SideEffects = append(c.SideEffects, EffectFetchRemote(remote))
	return nil
}

func (c *TestGitContext) ListAllBranches() ([]string, error) {
	return c.Branches, nil
}

func (c *TestGitContext) BranchExists(branchName string) (bool, error) {
	return slices.Contains(c.Branches, branchName), nil
}

func (c *TestGitContext) CheckoutBranch(branchName string) error {
	if !slices.Contains(c.Branches, branchName) {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	c.SideEffects = append(c.SideEffects, EffectCheckoutBranch(branchName))
	c.CurrentBranch = branchName
	return nil
}

func (c *TestGitContext) SwitchToNewBranch(branchName string) error {
	if slices.Contains(c.Branches, branchName) {
		return fmt.Errorf("branch %s already exists", branchName)
	}

	c.SideEffects = append(c.SideEffects, EffectCreateBranch(branchName))
	c.SideEffects = append(c.SideEffects, EffectCheckoutBranch(branchName))
	c.PreviousBranch = c.CurrentBranch
	c.CurrentBranch = branchName
	return nil
}

func (c *TestGitContext) SwitchBack() error {
	c.SideEffects = append(c.SideEffects, EffectSwitchBack())

	prev := c.PreviousBranch
	c.PreviousBranch = c.CurrentBranch
	c.CurrentBranch = prev
	return nil
}

func (c *TestGitContext) PushBranch(remote string, branchSpec string) error {
	c.SideEffects = append(c.SideEffects, EffectPushBranch(remote, branchSpec))
	return nil
}

func (c *TestGitContext) GetCurrentBranch() (string, error) {
	return c.CurrentBranch, nil
}

func (c *TestGitContext) ListRemotes() ([]string, error) {
	return c.Remotes, nil
}

func (c *TestGitContext) AssertNoSideEffects() {
	c.testCtx.Helper()
	if len(c.SideEffects) != 0 {
		c.testCtx.Fatalf("expected no side effects, got %s", formatSideEffects(c.SideEffects))
	}
}

func (c *TestGitContext) AssertSideEffectCount(expected int) {
	c.testCtx.Helper()
	if len(c.SideEffects) != expected {
		c.testCtx.Fatalf("expected %d side effects, got %d: %s", expected, len(c.SideEffects), formatSideEffects(c.SideEffects))
	}
}

func (c *TestGitContext) AssertSideEffectContains(expected TestGitSideEffect) {
	c.testCtx.Helper()
	if !slices.Contains(c.SideEffects, expected) {
		c.testCtx.Fatalf("expected side effects to contain '%s', but got %s", expected, formatSideEffects(c.SideEffects))
	}
}

func (c *TestGitContext) AssertSideEffectsAreExactly(expected ...TestGitSideEffect) {
	c.testCtx.Helper()

	if !slices.Equal(c.SideEffects, expected) {
		c.testCtx.Fatal(strings.Join([]string{
			"expected different side effects from what was executed",
			"expect: " + formatSideEffects(expected),
			"actual: " + formatSideEffects(c.SideEffects),
		}, "\n"))
	}
}

func formatSideEffects(sideEffects []TestGitSideEffect) string {
	formatted := []string{}
	for _, sideEffect := range sideEffects {
		formatted = append(formatted, fmt.Sprintf("'%s'", sideEffect))
	}
	return fmt.Sprintf("[%s]", strings.Join(formatted, ", "))
}

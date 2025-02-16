package gitrel_test

import (
	"fmt"
	"slices"
)

type TestGitContext struct {
	Branches       []string
	CurrentBranch  string
	PreviousBranch string
	Remotes        []string
	SideEffects    []TestGitSideEffect
}

func DefaultTestGitContext() *TestGitContext {
	return &TestGitContext{
		Branches: []string{
			"main",
			"release/1.0.0",
			"release/1.0.1",
			"release/1.0.2",
			"release/1.1.0",
			"release/1.1.1",
			"release/2.0.0",
			"remotes/origin/master",
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

	c.SideEffects = append(c.SideEffects, EffectSwitchToNewBranch(branchName))
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

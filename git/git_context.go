package git

import (
	"os/exec"
	"strings"
)

type GitContext interface {
	FetchRemote(remote string) error
	ListAllBranches() ([]string, error)
	BranchExists(branchName string) (bool, error)
	CheckoutBranch(branchName string) error
	SwitchToNewBranch(branchName string) error
	SwitchBack() error
	PushBranch(remote string, branchSpec string) error
	GetCurrentBranch() (string, error)
	ListRemotes() ([]string, error)
}

type CmdGitContext struct{}

func NewCmdGitContext() *CmdGitContext {
	return &CmdGitContext{}
}

func (c *CmdGitContext) FetchRemote(remote string) error {
	_, err := _execCommand("git", "fetch", remote)
	return err
}

func (c *CmdGitContext) ListAllBranches() ([]string, error) {
	output, err := _execCommand("git", "branch", "-a")
	if err != nil {
		return nil, err
	}

	rawBranches := strings.Split(output, "\n")
	branches := make([]string, 0, len(rawBranches))
	for _, branch := range rawBranches {
		if branch != "" {
			branches = append(branches, strings.TrimSpace(branch))
		}
	}

	return branches, nil
}

func (c *CmdGitContext) BranchExists(branchName string) (bool, error) {
	output, err := _execCommand("git", "branch", "--list", branchName)
	if err != nil {
		return false, err
	}

	return strings.Contains(output, branchName), nil
}

func (c *CmdGitContext) CheckoutBranch(branchName string) error {
	_, err := _execCommand("git", "checkout", branchName)
	return err
}

func (c *CmdGitContext) SwitchToNewBranch(branchName string) error {
	_, err := _execCommand("git", "switch", "-c", branchName)
	return err
}

func (c *CmdGitContext) SwitchBack() error {
	_, err := _execCommand("git", "switch", "-")
	return err
}

func (c *CmdGitContext) PushBranch(remote string, branchSpec string) error {
	_, err := _execCommand("git", "push", remote, branchSpec)
	return err
}

func (c *CmdGitContext) GetCurrentBranch() (string, error) {
	output, err := _execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

func (c *CmdGitContext) ListRemotes() ([]string, error) {
	output, err := _execCommand("git", "remote")
	if err != nil {
		return nil, err
	}

	rawRemotes := strings.Split(strings.TrimSpace(output), "\n")
	remotes := make([]string, 0, len(rawRemotes))
	for _, remote := range rawRemotes {
		if remote != "" {
			remotes = append(remotes, strings.TrimSpace(remote))
		}
	}

	return remotes, nil
}

// Function to execute a shell command and return its output
func _execCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

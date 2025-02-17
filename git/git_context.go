package git

import (
	"os/exec"
	"strings"
)

type CmdGitContext struct{}

func NewCmdGitContext() *CmdGitContext {
	return &CmdGitContext{}
}

func (c *CmdGitContext) HasUncommittedChanges() (bool, error) {
	output, err := _execCommand("git", "status", "--porcelain")
	if err != nil {
		return false, err
	}
	
	return output != "", nil
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
		branch = strings.TrimPrefix(branch, "*")
		
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

func (c *CmdGitContext) CreateBranchAt(branchName string, commitish string) error {
	_, err := _execCommand("git", "branch", branchName, commitish)
	return err
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

func (c *CmdGitContext) MergeBranch(branchName string) error {
	_, err := _execCommand("git", "merge", branchName)
	return err
}

// Function to execute a shell command and return its output
func _execCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

package interfaces

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
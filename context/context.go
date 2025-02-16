package context

type CommandContext interface {
	GetFetch() bool
	GetRemote() string
	GetLocalBranchName() string
	GetRemoteBranchName() string
}

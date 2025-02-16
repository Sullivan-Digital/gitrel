package interfaces

type CommandContext interface {
	GetOptFetch() bool
	GetOptRemote() string
	GetOptLocalBranchName() string
	GetOptRemoteBranchName() string

	SetFetched(fetched bool)
	GetFetched() bool
}

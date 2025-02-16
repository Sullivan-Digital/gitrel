package context

type CommandContext struct {
	Fetch            bool
	Remote           string
	LocalBranchName  string
	RemoteBranchName string
}
package gitrel_test

type TestCommandContext struct {
	Fetch            bool
	Remote           string
	LocalBranchName  string
	RemoteBranchName string
}

func DefaultTestCommandContext() *TestCommandContext {
	return &TestCommandContext{
		Fetch:            false,
		Remote:           "origin",
		LocalBranchName:  "release/%v",
		RemoteBranchName: "release/%v",
	}
}

func (c *TestCommandContext) GetFetch() bool {
	return c.Fetch
}

func (c *TestCommandContext) GetRemote() string {
	return c.Remote
}

func (c *TestCommandContext) GetLocalBranchName() string {
	return c.LocalBranchName
}

func (c *TestCommandContext) GetRemoteBranchName() string {
	return c.RemoteBranchName
}

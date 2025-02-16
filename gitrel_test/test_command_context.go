package gitrel_test

type TestCommandContext struct {
	Fetch            bool
	Remote           string
	LocalBranchName  string
	RemoteBranchName string

	fetched bool
}

func DefaultTestCommandContext() *TestCommandContext {
	return &TestCommandContext{
		Fetch:            false,
		Remote:           "origin",
		LocalBranchName:  "release/%v",
		RemoteBranchName: "release/%v",

		fetched: false,
	}
}

func (c *TestCommandContext) GetOptFetch() bool {
	return c.Fetch
}

func (c *TestCommandContext) GetOptRemote() string {
	return c.Remote
}

func (c *TestCommandContext) GetOptLocalBranchName() string {
	return c.LocalBranchName
}

func (c *TestCommandContext) GetOptRemoteBranchName() string {
	return c.RemoteBranchName
}

func (c *TestCommandContext) SetFetched(fetched bool) {
	c.fetched = fetched
}

func (c *TestCommandContext) GetFetched() bool {
	return c.fetched
}
package gitrel_test

import "gitrel/interfaces"

type TestGitRelContext struct {
	GitContext *TestGitContext
	CommandContext *TestCommandContext
	OutputContext *TestOutputContext
}

func DefaultTestGitRelContext() *TestGitRelContext {
	return &TestGitRelContext{
		GitContext: DefaultTestGitContext(),
		CommandContext: DefaultTestCommandContext(),
		OutputContext: DefaultTestOutputContext(),
	}
}

func (c *TestGitRelContext) Git() interfaces.GitContext {
	return c.GitContext
}

func (c *TestGitRelContext) Options() interfaces.CommandContext {
	return c.CommandContext
}

func (c *TestGitRelContext) Output() interfaces.OutputContext {
	return c.OutputContext
}

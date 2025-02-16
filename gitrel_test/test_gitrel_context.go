package gitrel_test

import (
	"gitrel/interfaces"
	"testing"
)

type TestGitRelContext struct {
	GitContext *TestGitContext
	CommandContext *TestCommandContext
	OutputContext *TestOutputContext
}

func DefaultTestGitRelContext(t *testing.T) *TestGitRelContext {
	return &TestGitRelContext{
		GitContext: DefaultTestGitContext(t),
		CommandContext: DefaultTestCommandContext(),
		OutputContext: DefaultTestOutputContext(t),
	}
}

func (c *TestGitRelContext) Git() interfaces.GitContext {
	return c.GitContext
}

func (c *TestGitRelContext) Command() interfaces.CommandContext {
	return c.CommandContext
}

func (c *TestGitRelContext) Output() interfaces.OutputContext {
	return c.OutputContext
}

package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"
)

type CmdGitRelContext struct {
	options interfaces.CommandContext
	git     interfaces.GitContext
	output  interfaces.OutputContext
}

func NewCmdGitRelContext() (*CmdGitRelContext, error) {
	gitCtx := git.NewCmdGitContext()
	cmdCtx, err := getCommandContext(gitCtx)
	if err != nil {
		return nil, err
	}

	ctx := &CmdGitRelContext{
		options: cmdCtx,
		git:     gitCtx,
		output:  NewCmdOutputContext(),
	}

	return ctx, nil
}

func (c *CmdGitRelContext) Options() interfaces.CommandContext {
	return c.options
}

func (c *CmdGitRelContext) Git() interfaces.GitContext {
	return c.git
}

func (c *CmdGitRelContext) Output() interfaces.OutputContext {
	return c.output
}

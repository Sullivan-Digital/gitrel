package cmd

import (
	"errors"
	"gitrel/config"
	"gitrel/context"
	"gitrel/git"
	"gitrel/utils"
)

var commandContext context.CommandContext

func getCommandContext() (context.CommandContext, error) {
	if commandContext != nil {
		return commandContext, nil
	}

	ctx := CmdCommandContext{}

	if FetchFlag && NoFetchFlag {
		return nil, errors.New("cannot use both --fetch and --no-fetch")
	}

	ctx.Fetch = FetchFlag || (config.FetchConfig && !NoFetchFlag)

	ctx.Remote = utils.CoalesceStr(RemoteFlag, config.RemoteConfig, "")
	if ctx.Remote == "" {
		remote, err := git.GetDefaultRemote()
		if err != nil {
			return nil, err
		}

		ctx.Remote = remote
	}

	ctx.LocalBranchName = utils.CoalesceStr(LocalBranchNameFlag, config.LocalBranchNameConfig, "release/%v")
	ctx.RemoteBranchName = utils.CoalesceStr(RemoteBranchNameFlag, config.RemoteBranchNameConfig, "release/%v")

	commandContext = &ctx
	return &ctx, nil
}

type CmdCommandContext struct {
	Fetch            bool
	Remote           string
	LocalBranchName  string
	RemoteBranchName string
}

func (c *CmdCommandContext) GetFetch() bool {
	return c.Fetch
}

func (c *CmdCommandContext) GetRemote() string {
	return c.Remote
}

func (c *CmdCommandContext) GetLocalBranchName() string {
	return c.LocalBranchName
}

func (c *CmdCommandContext) GetRemoteBranchName() string {
	return c.RemoteBranchName
}
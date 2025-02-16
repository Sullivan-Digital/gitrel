package cmd

import (
	"errors"
	"gitrel/config"
	"gitrel/git"
	"gitrel/utils"
	"gitrel/context"
)

var commandContext *context.CommandContext

func getCommandContext() (*context.CommandContext, error) {
	if commandContext != nil {
		return commandContext, nil
	}

	ctx := context.CommandContext{}

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
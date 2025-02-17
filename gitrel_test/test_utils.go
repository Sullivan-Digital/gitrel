package gitrel_test

import ()

func GetStdOutIgnoreSideEffects(ctx *TestGitRelContext, call func(*TestGitRelContext)) string {
	newCtx := &TestGitRelContext{
		GitContext: &TestGitContext{
			testCtx:        ctx.GitContext.testCtx,
			Branches:       ctx.GitContext.Branches,
			Remotes:        ctx.GitContext.Remotes,
			CurrentBranch:  ctx.GitContext.CurrentBranch,
			PreviousBranch: ctx.GitContext.PreviousBranch,
			SideEffects:    []TestGitSideEffect{},
		},
		CommandContext: &TestCommandContext{
			Fetch:            ctx.CommandContext.fetched,
			Remote:           ctx.CommandContext.Remote,
			LocalBranchName:  ctx.CommandContext.LocalBranchName,
			RemoteBranchName: ctx.CommandContext.RemoteBranchName,
			fetched:          ctx.CommandContext.fetched,
		},
		OutputContext: &TestOutputContext{
			Output: "",
			testCtx: ctx.OutputContext.testCtx,
		},
	}

	call(newCtx)
	return newCtx.OutputContext.Output
}

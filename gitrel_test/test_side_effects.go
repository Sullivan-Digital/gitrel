package gitrel_test

type TestGitSideEffect string

func EffectFetchRemote(remote string) TestGitSideEffect {
	return TestGitSideEffect("fetch " + remote)
}

func EffectCheckoutBranch(branch string) TestGitSideEffect {
	return TestGitSideEffect("checkout " + branch)
}

func EffectSwitchToNewBranch(branch string) TestGitSideEffect {
	return TestGitSideEffect("switch -c " + branch)
}

func EffectSwitchBack() TestGitSideEffect {
	return TestGitSideEffect("switch -")
}

func EffectPushBranch(remote string, branch string) TestGitSideEffect {
	return TestGitSideEffect("push " + remote + " " + branch)
}
package gitrel_test

import "strings"

type TestGitSideEffect string

func EffectFetchRemote(remote string) TestGitSideEffect {
	return TestGitSideEffect("fetch " + remote)
}

func EffectCreateBranch(branch string) TestGitSideEffect {
	return TestGitSideEffect("create branch " + branch)
}

func EffectCheckoutBranch(branch string) TestGitSideEffect {
	return TestGitSideEffect("checkout " + branch)
}

func EffectSwitchBack() TestGitSideEffect {
	return TestGitSideEffect("switch -")
}

func EffectPushBranch(remote string, branch string) TestGitSideEffect {
	parts := strings.Split(branch, ":")
	if len(parts) == 1 {
		return TestGitSideEffect("push " + remote + " " + branch)
	}

	if parts[0] == parts[1] {
		return TestGitSideEffect("push " + remote + " " + parts[0])
	}

	return TestGitSideEffect("push " + remote + " " + parts[0] + ":" + parts[1])
}
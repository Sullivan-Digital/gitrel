package git

type ReleaseInfo struct {
	Version  string
	Branches []ReleaseBranch
}

func (r *ReleaseInfo) IsLocalOnly() bool {
	for _, branch := range r.Branches {
		if branch.Type == "remote" {
			return false
		}
	}

	return true
}

func (r *ReleaseInfo) GetFirstLocalBranch() *ReleaseBranch {
	for _, branch := range r.Branches {
		if branch.Type == "local" {
			return &branch
		}
	}

	return nil
}

func (r *ReleaseInfo) GetFirstRemoteBranch() *ReleaseBranch {
	for _, branch := range r.Branches {
		if branch.Type == "remote" {
			return &branch
		}
	}

	return nil
}

type ReleaseBranch struct {
	BranchName string
	Type       string // remote or local
}

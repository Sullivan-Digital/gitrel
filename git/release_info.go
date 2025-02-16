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

type ReleaseBranch struct {
	Branch string
	Type   string // remote or local
}


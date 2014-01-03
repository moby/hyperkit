package vcs

func GetHgRepoRoot(path string) (isHgRepo bool, rootPath string) {
	// TODO: Not implemented
	return false, ""
}

type hgVcs struct {
	commonVcs
}

func (this *hgVcs) Type() Type { return Hg }

func (this *hgVcs) GetStatus() string {
	panic("TODO: Not implemented")
}

func (this *hgVcs) GetDefaultBranch() string {
	panic("TODO: Not implemented")
}

func (this *hgVcs) GetLocalBranch() string {
	panic("TODO: Not implemented")
}

func (this *hgVcs) GetLocalRev() string {
	panic("TODO: Not implemented")
}

func (this *hgVcs) GetRemoteRev() string {
	panic("TODO: Not implemented")
}

// Package for getting status of a repo under vcs.
package vcs

type Type uint8

const (
	Git Type = iota
	Hg
)

// TODO: Add comments.
type Vcs interface {
	RootPath() string
	Type() Type

	GetStatus() string
	GetDefaultBranch() string // TODO: Consider renaming GetRemoteBranch()
	GetLocalBranch() string
	GetLocalRev() string

	GetRemoteRev() string
}

type commonVcs struct {
	path string
}

func (this *commonVcs) RootPath() string {
	return this.path
}

// New returns a Vcs if path is under version control, otherwise nil.
func New(path string) Vcs {
	// TODO: This func should be done in a more general and smarter way
	if isRepo, rootPath := GetGitRepoRoot(path); isRepo {
		return &gitVcs{commonVcs{path: rootPath}}
	} else if isRepo, rootPath = GetHgRepoRoot(path); isRepo {
		return &hgVcs{commonVcs{path: rootPath}}
	} else {
		return nil
	}
}

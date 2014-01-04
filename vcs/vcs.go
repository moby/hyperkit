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

	GetDefaultBranch() string
	GetLocalBranch() string

	GetLocalRev() string
	GetRemoteRev() string
}

type commonVcs struct {
	rootPath string
}

func (this *commonVcs) RootPath() string {
	return this.rootPath
}

// New returns a Vcs if path is under version control, otherwise nil.
// TODO: Asking for same path should return point to existing Vcs, rather than creating another copy.
// Actually, maybe that should be the responsibility of a higher level package that uses this one, like VcsManager.
func New(path string) Vcs {
	for _, vcsProvider := range vcsProviders {
		if vcs := vcsProvider(path); vcs != nil {
			return vcs
		}
	}

	return nil
}

type vcsProvider func(path string) Vcs

var vcsProviders []vcsProvider

func addVcsProvider(s vcsProvider) {
	vcsProviders = append(vcsProviders, s)
}

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

// New returns a new Vcs if path is under version control, otherwise nil.
// It should be a valid path.
// TODO: Use a better type for path, e.g., github.com/shurcooL/go/path.
func New(path string) Vcs {
	// TODO: Optimize by checking more likely vcs first. Potentially check in parallel.
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

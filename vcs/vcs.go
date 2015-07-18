// Package vcs allows getting status of a repo under vcs.
package vcs

import "os/exec"

type Type uint8

const (
	Git Type = iota
	Hg
)

// VcsType returns a vcsType string compatible with sourcegraph.com/sourcegraph/go-vcs notation.
func (t Type) VcsType() (vcsType string) {
	switch t {
	case Git:
		return "git"
	case Hg:
		return "hg"
	default:
		panic("bad vcs.Type")
	}
}

type Vcs interface {
	RootPath() string // Returns the full path to the root of the repo.
	Type() Type       // Returns the type of vcs implementation.

	GetStatus() string // Returns empty string if no outstanding status.
	GetStash() string  // Returns empty string if no stash.

	GetRemote() string // Get primary remote repository url.

	GetDefaultBranch() string // Get default branch name for this vcs.
	GetLocalBranch() string   // Get currently checked out local branch name.

	GetLocalRev() string  // Get current local revision of default branch.
	GetRemoteRev() string // Get latest remote revision of default branch.

	// Returns true if given commit is contained in the default local branch.
	IsContained(rev string) bool
}

type commonVcs struct {
	rootPath string
}

func (this *commonVcs) RootPath() string {
	return this.rootPath
}

// New returns a new Vcs if path is under version control, otherwise nil.
// It should be a valid path.
func New(path string) Vcs {
	// TODO: Try to figure out vcs provider with a more constant-time operation.
	// TODO: Potentially check in parallel.
	for _, vcsProvider := range vcsProviders {
		if vcs := vcsProvider(path); vcs != nil {
			return vcs
		}
	}

	return nil
}

// Experimental, NewFromType returns a Vcs repository of the specified type without a local representation.
// Operations that require a local repository will fail.
func NewFromType(t Type) Vcs {
	switch t {
	case Git:
		return &gitVcs{}
	case Hg:
		return &hgVcs{}
	default:
		panic("bad vcs.Type")
	}
}

type vcsProvider func(path string) Vcs

var vcsProviders []vcsProvider

func addVcsProvider(s vcsProvider) {
	vcsProviders = append(vcsProviders, s)
}

func init() {
	// As an optimization, add Vcs providers sorted by the most likely first.

	// git.
	if _, err := exec.LookPath("git"); err == nil {
		addVcsProvider(func(path string) Vcs {
			if isRepo, rootPath := getGitRepoRoot(path); isRepo {
				return &gitVcs{commonVcs{rootPath: rootPath}}
			}
			return nil
		})
	}

	// hg.
	if _, err := exec.LookPath("hg"); err == nil {
		addVcsProvider(func(path string) Vcs {
			if isRepo, rootPath := getHgRepoRoot(path); isRepo {
				return &hgVcs{commonVcs{rootPath: rootPath}}
			}
			return nil
		})
	}
}

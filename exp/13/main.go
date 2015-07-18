// Package exp13 offers caching of vcs state per repository.
package exp13

import (
	"github.com/shurcooL/go/gists/gist7802150"
	"github.com/shurcooL/go/vcs"
	go_vcs "golang.org/x/tools/go/vcs"
)

type VcsState struct {
	Vcs vcs.Vcs

	VcsLocal  *VcsLocal
	VcsRemote *VcsRemote

	RepoRoot *go_vcs.RepoRoot

	// THINK: No need to add repo as a DepNode2I, just add it a plain variable. Maybe?
	// TODO: No need for this to have a DepNode2Manual, remove it.
	//       Well, the idea is I don't foresee anyone invalidating the entire VcsState.
	gist7802150.DepNode2Manual
}

func NewVcsState(vcs vcs.Vcs) *VcsState {
	this := &VcsState{
		Vcs: vcs,
	}
	this.VcsLocal = NewVcsLocal(this)
	this.VcsRemote = NewVcsRemote(this)
	return this
}

// ---

type VcsLocal struct {
	Status      string
	Stash       string
	Remote      string
	LocalBranch string
	LocalRev    string

	gist7802150.DepNode2
}

func NewVcsLocal(repo *VcsState) *VcsLocal {
	this := &VcsLocal{}
	// THINK: No need to add repo as a DepNode2I, just add it a plain variable. Maybe?
	this.AddSources(repo, &gist7802150.DepNode2Manual{})
	return this
}

func (this *VcsLocal) Update() {
	// THINK: No need to add repo as a DepNode2I, just add it a plain variable. Maybe?
	vcs := this.GetSources()[0].(*VcsState).Vcs

	this.Status = vcs.GetStatus()
	this.Stash = vcs.GetStash()
	this.Remote = vcs.GetRemote()
	this.LocalBranch = vcs.GetLocalBranch()
	this.LocalRev = vcs.GetLocalRev()
}

// ---

type VcsRemote struct {
	RemoteRev   string
	IsContained bool // True if remote commit is contained in the default local branch.

	gist7802150.DepNode2
}

func NewVcsRemote(repo *VcsState) *VcsRemote {
	this := &VcsRemote{}
	this.AddSources(repo)
	return this
}

func (this *VcsRemote) Update() {
	vcs := this.GetSources()[0].(*VcsState).Vcs

	this.RemoteRev = vcs.GetRemoteRev()
	if this.RemoteRev != "" {
		this.IsContained = vcs.IsContained(this.RemoteRev)
	}
}

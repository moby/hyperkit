package exp13

import (
	. "github.com/shurcooL/go/gists/gist7802150"

	"github.com/shurcooL/go/vcs"
)

type VcsState struct {
	Vcs vcs.Vcs

	VcsLocal  *VcsLocal
	VcsRemote *VcsRemote

	// THINK: No need to add repo as a DepNode2I, just add it a plain variable. Maybe?
	// TODO: No need for this to have a DepNode2Manual, remove it.
	//       Well, the idea is I don't foresee anyone invalidating the entire VcsState.
	DepNode2Manual
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

	DepNode2
}

func NewVcsLocal(repo *VcsState) *VcsLocal {
	this := &VcsLocal{}
	// THINK: No need to add repo as a DepNode2I, just add it a plain variable. Maybe?
	this.AddSources(repo, &DepNode2Manual{})
	return this
}

func (this *VcsLocal) Update() {
	// THINK: No need to add repo as a DepNode2I, just add it a plain variable. Maybe?
	vcs := this.GetSources()[0].(*VcsState).Vcs

	//fmt.Println("*VcsLocal) Update() for", vcs.RootPath())

	this.Status = vcs.GetStatus()
	this.Stash = vcs.GetStash()
	this.Remote = vcs.GetRemote()
	this.LocalBranch = vcs.GetLocalBranch()
	this.LocalRev = vcs.GetLocalRev()
}

// ---

type VcsRemote struct {
	RemoteRev string

	DepNode2
}

func NewVcsRemote(repo *VcsState) *VcsRemote {
	this := &VcsRemote{}
	this.AddSources(repo)
	return this
}

func (this *VcsRemote) Update() {
	vcs := this.GetSources()[0].(*VcsState).Vcs

	//fmt.Println("*VcsRemote) Update() for", vcs.RootPath())

	this.RemoteRev = vcs.GetRemoteRev()
}

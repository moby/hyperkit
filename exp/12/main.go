package exp12

import (
	"github.com/shurcooL/go/exp/13"
	"github.com/shurcooL/go/vcs"

	. "gist.github.com/7802150.git"
)

// TODO: Rename to "Folder" or "FileSystemNode" or something.
type MaybeVcsRepo struct {
	path string

	VcsState *exp13.VcsState

	DepNode2
}

func (this *MaybeVcsRepo) Update() {
	//this.Vcs = vcs.New(this.path)

	if vcs := vcs.New(this.path); vcs != nil {
		if vcsState, ok := vcsStates[vcs.RootPath()]; ok {
			this.VcsState = vcsState
		} else {
			this.VcsState = exp13.NewVcsState(vcs)
			vcsStates[vcs.RootPath()] = this.VcsState
		}
	}
}

func NewMaybeVcsRepo(path string) *MaybeVcsRepo {
	this := &MaybeVcsRepo{path: path}
	// No DepNode2I sources, so each instance can only be updated (i.e. initialized) once
	return this
}

// =====

// TODO: Use FileUri or similar instead of string for clean path to repo root.
// rootPath -> *VcsState
var vcsStates = make(map[string]*exp13.VcsState)

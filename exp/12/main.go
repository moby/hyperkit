package exp12

import (
	"github.com/shurcooL/go/exp/13"
	"github.com/shurcooL/go/vcs"

	. "gist.github.com/7802150.git"
)

// TODO: Use FileUri or similar type instead of string for clean path to repo root.
// rootPath -> *VcsState
var repos = make(map[string]*exp13.VcsState)

type Directory struct {
	path string

	Repo *exp13.VcsState

	DepNode2
}

func (this *Directory) Update() {
	if vcs := vcs.New(this.path); vcs != nil {
		if repo, ok := repos[vcs.RootPath()]; ok {
			this.Repo = repo
		} else {
			this.Repo = exp13.NewVcsState(vcs)
			repos[vcs.RootPath()] = this.Repo
		}
	}
}

func NewDirectory(path string) *Directory {
	this := &Directory{path: path}
	// No DepNode2I sources, so each instance can only be updated (i.e. initialized) once
	return this
}

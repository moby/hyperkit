package exp14

import (
	. "github.com/shurcooL/go/gists/gist7480523"
	. "github.com/shurcooL/go/gists/gist7802150"

	"github.com/shurcooL/go/gists/gist8018045"
)

type GoPackageList interface {
	List() []*GoPackage

	DepNode2I
}

type GoPackages struct {
	SkipGoroot bool // Currently, works on initial run only; changing its value afterwards has no effect.

	Entries []*GoPackage

	DepNode2
}

func (this *GoPackages) Update() {
	// TODO: Have a source?

	// TODO: Make it load in background, without blocking, etc.
	{
		goPackages := make(chan *GoPackage, 64)

		if this.SkipGoroot {
			go gist8018045.GetGopathGoPackages(goPackages)
		} else {
			go gist8018045.GetGoPackages(goPackages)
		}

		this.Entries = nil
		for {
			if goPackage, ok := <-goPackages; ok {
				this.Entries = append(this.Entries, goPackage)
			} else {
				break
			}
		}
	}
}

func (this *GoPackages) List() []*GoPackage {
	return this.Entries
}

package exp14

import (
	"io"

	. "github.com/shurcooL/go/gists/gist7480523"
	"github.com/shurcooL/go/gists/gist7651991"
	. "github.com/shurcooL/go/gists/gist7802150"

	"github.com/shurcooL/go/gists/gist8018045"
)

type GoPackageList interface {
	List() []*GoPackage

	DepNode2I
}

// GoPackages is a cached list of all Go packages in GOPATH including/excluding GOROOT.
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

// GoPackagesFromReader is a cached list of Go packages specified by newline separated import paths from Reader.
type GoPackagesFromReader struct {
	Reader io.Reader

	Entries []*GoPackage

	DepNode2
}

func (this *GoPackagesFromReader) Update() {
	reduceFunc := func(importPath string) interface{} {
		if goPackage := GoPackageFromImportPath(importPath); goPackage != nil {
			return goPackage
		}
		return nil
	}

	goPackages := gist7651991.GoReduceLinesFromReader(this.Reader, 8, reduceFunc)

	this.Entries = nil
	for {
		if goPackage, ok := <-goPackages; ok {
			this.Entries = append(this.Entries, goPackage.(*GoPackage))
		} else {
			break
		}
	}
}

func (this *GoPackagesFromReader) List() []*GoPackage {
	return this.Entries
}

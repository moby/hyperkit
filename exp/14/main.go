package exp14

import (
	. "gist.github.com/7480523.git"
	. "gist.github.com/7802150.git"

	"gist.github.com/8018045.git"
)

type GoPackages struct {
	Entries []*GoPackage

	DepNode2
}

func (this *GoPackages) Update() {
	// TODO: Have a source?

	// TODO: Make it load in background, without blocking, etc.
	{
		goPackages := make(chan *GoPackage, 64)

		go gist8018045.GetGoPackages(goPackages)

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

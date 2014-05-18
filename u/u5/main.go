// Package u5 currently provides a single utility to fetch the importers of a GoPackage via godoc.org API.
package u5

import (
	"encoding/json"
	"net/http"

	"gist.github.com/7480523.git"
)

type goPackage struct {
	Path     string
	Synopsis string
}

// Importers contains the list of Go packages that import a given Go package.
type Importers struct {
	Results []goPackage
}

// GetGodocOrgImporters fetches the importers of goPackage via godoc.org API.
func GetGodocOrgImporters(goPackage *gist7480523.GoPackage) (*Importers, error) {
	resp, err := http.Get("http://api.godoc.org/importers/" + goPackage.Bpkg.ImportPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var importers Importers
	if err := json.NewDecoder(resp.Body).Decode(&importers); err != nil {
		return nil, err
	}

	return &importers, nil
}

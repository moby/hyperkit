// Package u5 currently provides a single utility to fetch the importers of a GoPackage via godoc.org API.
package u5

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shurcooL/go/gists/gist7480523"
)

// GoPackage represents a Go package.
type GoPackage struct {
	Path     string // Import path of the package.
	Synopsis string // Synopsis of the package.
}

// Importers contains the list of Go packages that import a given Go package.
type Importers struct {
	Results []GoPackage
}

// GetGodocOrgImporters fetches the importers of goPackage via godoc.org API.
func GetGodocOrgImporters(goPackage *gist7480523.GoPackage) (*Importers, error) {
	resp, err := http.Get("http://api.godoc.org/importers/" + goPackage.Bpkg.ImportPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %v", resp.StatusCode)
	}

	var importers Importers
	if err := json.NewDecoder(resp.Body).Decode(&importers); err != nil {
		return nil, err
	}

	return &importers, nil
}

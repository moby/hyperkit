// Package u5 provides a single utility to fetch the importers of a GoPackage via godoc.org API.
package u5

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// UserAgent is used for outbound requests to godoc.org API, if set to non-empty value.
var UserAgent string

// GetGodocOrgImporters fetches the importers of Go package with specified importPath via godoc.org API.
func GetGodocOrgImporters(importPath string) (*Importers, error) {
	req, err := http.NewRequest("GET", "http://api.godoc.org/importers/"+importPath, nil)
	if err != nil {
		return nil, err
	}
	if UserAgent != "" {
		req.Header.Set("User-Agent", UserAgent)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %v", resp.StatusCode)
	}
	var importers Importers
	err = json.NewDecoder(resp.Body).Decode(&importers)
	if err != nil {
		return nil, err
	}
	return &importers, nil
}

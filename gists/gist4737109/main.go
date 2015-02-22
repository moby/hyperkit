// Package gist4737109 gets the contents of a gist.
package gist4737109

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GistIdToUsername returns the GitHub username owner of gist with given gistId.
func GistIdToUsername(gistId string) (string, error) {
	gistUrl := "https://api.github.com/gists/" + gistId

	gistBytes, err := httpGet(gistUrl)
	if err != nil {
		return "", err
	}
	var gistJson struct {
		Owner struct{ Login string }
	}
	err = json.Unmarshal(gistBytes, &gistJson)
	if err != nil {
		return "", err
	}
	return gistJson.Owner.Login, nil
}

func main() {
	gistId := "4737109"

	fmt.Println(GistIdToUsername(gistId))
}

// ---

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %v", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

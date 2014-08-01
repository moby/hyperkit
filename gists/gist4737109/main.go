// Package gist4737109 gets the contents of a gist.
package gist4737109

import (
	"encoding/json"
	"fmt"

	. "github.com/shurcooL/go/gists/gist4668739"
)

/* TODO: Probably want to change interface to return error, etc. But it's worth doing when these funcs are used, not sooner.
// This assumes there's only one file in the gist
func GistIdToGistContents(gistId string) string {
	return GistIdCommitIdToGistContents(gistId, "")
}

// This assumes there's only one file in the gist
func GistIdCommitIdToGistContents(gistId, commitId string) string {
	gistUrl := "https://api.github.com/gists/" + gistId
	if "" != commitId {
		gistUrl += "/" + commitId
	}

	var gistJson struct {
		Files map[string]struct{ Content string }
		//History []struct{ Version string }
	}
	err := json.Unmarshal(HttpGetB(gistUrl), &gistJson)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, v := range gistJson.Files {
		return v.Content
	}
	return ""
}*/

func GistIdToUsername(gistId string) (string, error) {
	gistUrl := "https://api.github.com/gists/" + gistId

	var gistJson struct {
		Owner struct{ Login string }
	}
	err := json.Unmarshal(HttpGetB(gistUrl), &gistJson)
	if err != nil {
		return "", err
	}
	return gistJson.Owner.Login, nil
}

func main() {
	gistId := "4737109"
	//commitId := "1b4e4b0e6f469d5e5a91b49028fbf2ab936bfcd4"

	fmt.Println(GistIdToUsername(gistId))
	//println(GistIdCommitIdToGistContents(gistId, commitId))
	//println(GistIdToGistContents(gistId))
}

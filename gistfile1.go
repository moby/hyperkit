package gist4737109

import (
	"encoding/json"
	"fmt"
	. "gist.github.com/4668739.git"
	//"strconv"
)

// This assumes there's only one file in the gist
func GistIdToGistContents(gistId string) string {
	return GistIdCommitIdToGistContents(gistId, "")
}

// This assumes there's only one file in the gist
func GistIdCommitIdToGistContents(gistId, commitId string) string {
	var out string

	gistUrl := "https://api.github.com/gists/" + gistId
	if ("" != commitId) {
		gistUrl = "https://api.github.com/gists/" + gistId + "/" + commitId
	}
	b := HttpGetB(gistUrl)

	type GistFile struct {
		Raw_Url string
	}
	type Response struct {
		Files map[string]GistFile
	}

	var animals Response
	err := json.Unmarshal(b, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, v := range animals.Files {
		out = v.Raw_Url
		break
	}
	return HttpGet(out)
}

func main() {
	gistId := "4737109"
	commitId := "1b4e4b0e6f469d5e5a91b49028fbf2ab936bfcd4"

	println(GistIdCommitIdToGistContents(gistId, commitId))
	println(GistIdToGistContents(gistId))
}
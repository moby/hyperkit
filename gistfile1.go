package gist4737109

import (
	"encoding/json"
	"fmt"
	. "gist.github.com/4668739.git"
	//"strconv"
)

// This assumes there's only one file in the gist
func GistIdToGistContents(gistId string) string {
	var out string

	b := HttpGetB("https://api.github.com/gists/" + gistId)

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

	print(GistIdToGistContents(gistId))
}
// Package gist4737109 gets the contents of a gist.
package gist4737109

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GistIdToUsername(gistId string) (string, error) {
	gistUrl := "https://api.github.com/gists/" + gistId

	var gistJson struct {
		Owner struct{ Login string }
	}
	err := json.Unmarshal(httpGetB(gistUrl), &gistJson)
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

func httpGetB(url string) []byte {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return b
}

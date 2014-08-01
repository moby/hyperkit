// Package gist4668739 gets the contents of a webpage as a string.
package gist4668739

import (
	"io/ioutil"
	"net/http"

	. "github.com/shurcooL/go/gists/gist5286084"
)

func HttpGet(url string) string {
	return string(HttpGetB(url))
}

func HttpGetB(url string) []byte {
	r, err := http.Get(url)
	CheckError(err)
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	CheckError(err)
	return b
}

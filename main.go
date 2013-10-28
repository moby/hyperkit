package gist4668739

import (
	"io/ioutil"
	"net/http"
	. "gist.github.com/5286084.git"
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
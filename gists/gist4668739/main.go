// Package gist4668739 gets the contents of a webpage as a string.
package gist4668739

import (
	"io/ioutil"
	"net/http"
)

// DEPRECATED.
func HttpGet(url string) string {
	return string(HttpGetB(url))
}

// DEPRECATED.
func HttpGetB(url string) []byte {
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

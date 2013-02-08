package gist4668739

import (
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) string {
	return string(HttpGetB(url))
}

func HttpGetB(url string) []byte {
	r, err := http.Get(url)
	if nil != err {
		panic(err)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if nil != err {
		panic(err)
	}
	return b
}
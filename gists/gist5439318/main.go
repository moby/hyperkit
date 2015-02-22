// Package gist5439318 gets the contents of a tweet.
package gist5439318

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/shurcooL/go-goon"
)

func GetTweet(id string) map[string]interface{} {
	tweetBytes := httpGetB("https://api.twitter.com/1/statuses/oembed.json?id=" + id + "&omit_script=true")
	var tweetJson map[string]interface{}
	err := json.Unmarshal(tweetBytes, &tweetJson)
	if err != nil {
		panic(err)
	}
	return tweetJson
}

func GetTweetHtml(id string) string {
	t := GetTweet(id)
	return t["html"].(string)
}

func main() {
	goon.Dump(GetTweet("289608996225171456"))
	goon.Dump(GetTweetHtml("289608996225171456"))
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

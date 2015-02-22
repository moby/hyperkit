// Package gist5439318 gets the contents of a tweet.
package gist5439318

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/shurcooL/go-goon"
)

func GetTweet(id string) map[string]interface{} {
	tweetBytes, err := httpGet("https://api.twitter.com/1/statuses/oembed.json?id=" + id + "&omit_script=true")
	if err != nil {
		panic(err)
	}
	var tweetJson map[string]interface{}
	err = json.Unmarshal(tweetBytes, &tweetJson)
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

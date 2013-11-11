package main

import (
	"encoding/json"
	. "gist.github.com/4668739.git"
	. "gist.github.com/5286084.git"

	"github.com/shurcooL/go-goon"
)

func GetTweet(id string) map[string]interface{} {
	tweetBytes := HttpGetB("https://api.twitter.com/1/statuses/oembed.json?id=" + id + "&omit_script=true")
	var tweetJson map[string]interface{}
	err := json.Unmarshal(tweetBytes, &tweetJson)
	CheckError(err)
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

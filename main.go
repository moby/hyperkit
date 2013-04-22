package main

import (
	. "gist.github.com/4668739.git"
	. "gist.github.com/5286084.git"
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
)

func GetTweet(id string) map[string]interface{} {
	jsonb := HttpGetB("https://api.twitter.com/1/statuses/oembed.json?id=" + id + "&omit_script=true")
	var f map[string]interface{}
	err := json.Unmarshal(jsonb, &f)
	CheckError(err)
	return f
}

func GetTweetHtml(id string) string {
	t := GetTweet(id)
	return t["html"].(string)
}

func main() {
	spew.Dump(GetTweet("289608996225171456"))
	spew.Dump(GetTweetHtml("289608996225171456"))
}
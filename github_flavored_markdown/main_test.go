package github_flavored_markdown_test

import (
	"os"

	"github.com/shurcooL/go/github_flavored_markdown"
)

func ExampleMarkdown() {
	text := []byte("Hello world github/linguist#1 **cool**, and #1!")

	os.Stdout.Write(github_flavored_markdown.Markdown(text))

	// Output:
	//<p>Hello world github/linguist#1 <strong>cool</strong>, and #1!</p>
}

func ExampleHeader() {
	text := []byte("## git diff")

	os.Stdout.Write(github_flavored_markdown.Markdown(text))

	// Output:
	//<h2><a name="git-diff" class="anchor" href="#git-diff" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>git diff</h2>
}

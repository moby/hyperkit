package u1

import (
	"bytes"
	"io"
	"net/http"

	"github.com/shurcooL/go/github_flavored_markdown"
)

// Convert GitHub Flavored Markdown to full HTML page and write to w.
// TODO: Do this locally via a native Go library... That's not too much to ask for, is it?
func WriteMarkdownGfmAsHtmlPage(w io.Writer, markdown []byte) {
	// TODO: Do GitHub, fallback to local if it fails.
	writeGitHubFlavoredMarkdownViaGitHub(w, markdown)
	//WriteGitHubFlavoredMarkdownViaLocal(w, markdown)
}

func WriteGitHubFlavoredMarkdownViaLocal(w io.Writer, markdown []byte) {
	// TODO: Don't hotlink the css file from github.com, serve it locally (it's needed for the GFM html to appear properly)
	io.WriteString(w, `<html><head><meta charset="utf-8"><style>code, div.highlight { tab-size: 4; }</style><link href="https://github.com/assets/github.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(github_flavored_markdown.Markdown(markdown))
	io.WriteString(w, `</article></body></html>`)
}

// TODO: Remove once local version gives matching results.
func writeGitHubFlavoredMarkdownViaGitHub(w io.Writer, markdown []byte) {
	// TODO: Don't hotlink the css file from github.com, serve it locally (it's needed for the GFM html to appear properly)
	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="https://github.com/assets/github.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)

	// Convert GitHub-Flavored-Markdown to HTML (includes syntax highlighting for diff, Go, etc.)
	resp, err := http.Post("https://api.github.com/markdown/raw", "text/x-markdown", bytes.NewReader(markdown))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}

	io.WriteString(w, `</article></body></html>`)
}

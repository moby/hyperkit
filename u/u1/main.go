package u1

import (
	"bytes"
	"io"
	"net/http"

	"github.com/shurcooL/go/github_flavored_markdown"
)

// Convert GitHub Flavored Markdown to full HTML page and write to w.
func WriteGitHubFlavoredMarkdownViaLocal(w io.Writer, markdown []byte) {
	// TODO: Don't hotlink the css file from github.com, serve it locally (it's needed for the GFM html to appear properly)
	io.WriteString(w, `<html><head><meta charset="utf-8"><style>code, div.highlight { tab-size: 4; }</style><link href="https://assets-cdn.github.com/assets/github-aac267329edcf91101f8b8afc0575e8e882615cc.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(github_flavored_markdown.Markdown(markdown))
	io.WriteString(w, `</article></body></html>`)
}

// TODO: Remove once local version gives matching results.
func WriteGitHubFlavoredMarkdownViaGitHub(w io.Writer, markdown []byte) {
	// TODO: Don't hotlink the css file from github.com, serve it locally (it's needed for the GFM html to appear properly)
	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="https://assets-cdn.github.com/assets/github-aac267329edcf91101f8b8afc0575e8e882615cc.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)

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

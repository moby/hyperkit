package u1

import (
	"bytes"
	"io"
	"net/http"

	"github.com/russross/blackfriday"
)

// GitHub Flavored Markdown-like extensions.
var MarkdownGfmExtensions = 0 |
	blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	//blackfriday.EXTENSION_TABLES | // TODO: Implement. Maybe.
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS
	//blackfriday.EXTENSION_HARD_LINE_BREAK

// Best effort at generating GitHub Flavored Markdown-like HTML output locally.
func MarkdownGfm(input []byte) []byte {
	htmlFlags := 0 |
		blackfriday.HTML_SANITIZE_OUTPUT |
		blackfriday.HTML_GITHUB_BLOCKCODE

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	return blackfriday.Markdown(input, renderer, MarkdownGfmExtensions)
}

// ---

// Convert GitHub Flavored Markdown to full HTML page and write to w.
// TODO: Do this locally via a native Go library... That's not too much to ask for, is it?
func WriteMarkdownGfmAsHtmlPage(w io.Writer, markdown []byte) {
	// TODO: Do GitHub, fallback to local if it fails.
	writeGitHubFlavoredMarkdownViaGitHub(w, markdown)
	//WriteGitHubFlavoredMarkdownViaLocal(w, markdown)
}

func WriteGitHubFlavoredMarkdownViaLocal(w io.Writer, markdown []byte) {
	// TODO: Don't hotlink the css file from github.com, serve it locally (it's needed for the GFM html to appear properly)
	// TODO: Use github.com/sourcegraph/syntaxhighlight to add missing syntax highlighting.
	io.WriteString(w, `<html><head><meta charset="utf-8"><style>code, div.highlight { tab-size: 4; }</style><link href="https://github.com/assets/github.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(MarkdownGfm(markdown))
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

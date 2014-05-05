package u1

import "github.com/russross/blackfriday"

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

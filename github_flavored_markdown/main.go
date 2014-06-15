// Package github_flavored_markdown provides a GitHub Flavored Markdown renderer
// with fenced code block highlighting.
package github_flavored_markdown

import (
	"bytes"
	"strings"

	"github.com/russross/blackfriday"
	"github.com/shurcooL/go/u/u7"
	"github.com/sourcegraph/syntaxhighlight"
)

// Markdown is a convenience function for rendering input GitHub Flavored Markdown.
//
// It does not attempt to sanitize HTML output; you can do that in post-processing using github.com/microcosm-cc/bluemonday package.
func Markdown(input []byte) []byte {
	renderer := NewRenderer()

	// GitHub Flavored Markdown extensions.
	var gfmExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS
		//blackfriday.EXTENSION_HARD_LINE_BREAK

	return blackfriday.Markdown(input, renderer, gfmExtensions)
}

// NewRenderer creates a GitHub Flavored Markdown HTML renderer, which satisfies the blackfriday.Renderer interface.
//
// It does not attempt to sanitize HTML output; you can do that in post-processing using github.com/microcosm-cc/bluemonday package.
func NewRenderer() blackfriday.Renderer {
	htmlFlags := 0 |
		//blackfriday.HTML_SANITIZE_OUTPUT |
		blackfriday.HTML_GITHUB_BLOCKCODE

	return &renderer{blackfriday.HtmlRenderer(htmlFlags, "", "").(*blackfriday.Html)}
}

type renderer struct {
	*blackfriday.Html
}

// TODO: Clean up and improve this code.
// GitHub Flavored Markdown fenced code block with highlighting.
func (r *renderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	doubleSpace(out)

	// parse out the language name
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		out.WriteString(`<div class="highlight highlight-`)
		attrEscape(out, []byte(elt))
		lang = elt
		out.WriteString(`"><pre>`)
		count++
		break
	}

	if count == 0 {
		out.WriteString("<pre><code>")
	}

	if formattedCode, ok := formatCode(text, lang); ok {
		out.Write(formattedCode)
	} else {
		attrEscape(out, text)
	}

	if count == 0 {
		out.WriteString("</code></pre>\n")
	} else {
		out.WriteString("</pre></div>\n")
	}
}

var gfmHTMLConfig = syntaxhighlight.HTMLConfig{
	String:        "s",
	Keyword:       "k",
	Comment:       "c",
	Type:          "n",
	Literal:       "lit",
	Punctuation:   "p",
	Plaintext:     "n",
	Tag:           "tag",
	HTMLTag:       "htm",
	HTMLAttrName:  "atn",
	HTMLAttrValue: "atv",
	Decimal:       "m",
}

// TODO: Support highlighting for more languages.
func formatCode(src []byte, lang string) (formattedCode []byte, ok bool) {
	switch lang {
	// TODO: Use a highlighter based on go/scanner for Go code.
	case "Go", "go":
		var buf bytes.Buffer
		err := syntaxhighlight.Print(syntaxhighlight.NewScanner(src), &buf, syntaxhighlight.HTMLPrinter(gfmHTMLConfig))
		if err != nil {
			return nil, false
		}
		return buf.Bytes(), true
	case "diff":
		var buf bytes.Buffer
		err := u7.Print(u7.NewScanner(src), &buf)
		if err != nil {
			return nil, false
		}
		return buf.Bytes(), true
	default:
		return nil, false
	}
}

// Unexported blackfriday helpers.

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

func escapeSingleChar(char byte) (string, bool) {
	if char == '"' {
		return "&quot;", true
	}
	if char == '&' {
		return "&amp;", true
	}
	if char == '<' {
		return "&lt;", true
	}
	if char == '>' {
		return "&gt;", true
	}
	return "", false
}

func attrEscape(out *bytes.Buffer, src []byte) {
	org := 0
	for i, ch := range src {
		if entity, ok := escapeSingleChar(ch); ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}
}

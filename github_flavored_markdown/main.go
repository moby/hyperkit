/*
Package github_flavored_markdown provides a GitHub Flavored Markdown renderer
with fenced code block highlighting, clickable header anchor links.

The functionality should be equivalent to the GitHub Markdown API endpoint specified at
https://developer.github.com/v3/markdown/#render-a-markdown-document-in-raw-mode, except
the rendering is performed locally.

See example below for how to generate a complete HTML page, including CSS styles.
*/
package github_flavored_markdown

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/shurcooL/go/u/u7"
	"github.com/sourcegraph/annotate"
	"github.com/sourcegraph/syntaxhighlight"
)

// Markdown renders GitHub Flavored Markdown text.
func Markdown(text []byte) []byte {
	htmlFlags := 0
	//htmlFlags |= blackfriday.HTML_SANITIZE_OUTPUT
	htmlFlags |= blackfriday.HTML_GITHUB_BLOCKCODE
	renderer := &renderer{Html: blackfriday.HtmlRenderer(htmlFlags, "", "").(*blackfriday.Html)}

	// Parser extensions for GitHub Flavored Markdown.
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	//extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK

	unsanitized := blackfriday.Markdown(text, renderer, extensions)

	// GitHub Flavored Markdown-like sanitization policy.
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("div", "span")
	p.AllowAttrs("class", "name").Matching(bluemonday.SpaceSeparatedTokens).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^nofollow$`)).OnElements("a")
	p.AllowAttrs("aria-hidden").Matching(regexp.MustCompile(`^true$`)).OnElements("a")
	p.AllowDataURIImages()

	return p.SanitizeBytes(unsanitized)
}

type renderer struct {
	*blackfriday.Html
}

// GitHub Flavored Markdown header with clickable and hidden anchor.
func (_ *renderer) Header(out *bytes.Buffer, text func() bool, level int, _ string) {
	marker := out.Len()
	doubleSpace(out)

	if !text() {
		out.Truncate(marker)
		return
	}

	textString := out.String()[marker:]
	out.Truncate(marker)

	anchorName := createSanitizedAnchorName(html.UnescapeString(textString))

	out.WriteString(fmt.Sprintf(`<h%d><a name="%s" class="anchor" href="#%s" rel="nofollow" aria-hidden="true"><span class="octicon octicon-link"></span></a>`, level, anchorName, anchorName))
	out.WriteString(textString)
	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

// Returns an anchor name for the given header text.
func createSanitizedAnchorName(text string) string {
	var anchorName []rune
	for _, r := range []rune(text) {
		switch {
		case r == ' ':
			anchorName = append(anchorName, '-')
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			anchorName = append(anchorName, unicode.ToLower(r))
		}
	}
	return string(anchorName)
}

// TODO: Clean up and improve this code.
// GitHub Flavored Markdown fenced code block with highlighting.
func (_ *renderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
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

var gfmHtmlConfig = syntaxhighlight.HTMLConfig{
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
		err := syntaxhighlight.Print(syntaxhighlight.NewScanner(src), &buf, syntaxhighlight.HTMLPrinter(gfmHtmlConfig))
		if err != nil {
			return nil, false
		}
		return buf.Bytes(), true
	case "diff":
		switch 2 {
		default:
			var buf bytes.Buffer
			err := u7.Print(u7.NewScanner(src), &buf)
			if err != nil {
				return nil, false
			}
			return buf.Bytes(), true
		case 1:
			lines := bytes.Split(src, []byte("\n"))
			return bytes.Join(lines, []byte("\n")), true
		case 2:
			anns, err := u7.Annotate(src)
			if err != nil {
				return nil, false
			}

			lines := bytes.Split(src, []byte("\n"))
			lineStarts := make([]int, len(lines))
			var offset int
			for lineIndex := 0; lineIndex < len(lines); lineIndex++ {
				lineStarts[lineIndex] = offset
				offset += len(lines[lineIndex]) + 1
			}

			lastDel, lastIns := -1, -1
			for lineIndex := 0; lineIndex < len(lines); lineIndex++ {
				var lineFirstChar byte
				if len(lines[lineIndex]) > 0 {
					lineFirstChar = lines[lineIndex][0]
				}
				switch lineFirstChar {
				case '+':
					if lastIns == -1 {
						lastIns = lineIndex
					}
				case '-':
					if lastDel == -1 {
						lastDel = lineIndex
					}
				default:
					if lastDel != -1 && lastIns != -1 && lastDel+1 == lastIns && lastIns+1 == lineIndex && '@' != lineFirstChar {
						if lastDel == -1 {
							lastDel = lastIns
						} else if lastIns == -1 {
							lastIns = lineIndex
						}

						beginOffsetLeft := lineStarts[lastDel] + 1
						endOffsetLeft := lineStarts[lastIns]
						beginOffsetRight := lineStarts[lastIns] + 1
						endOffsetRight := lineStarts[lineIndex]

						leftContent := string(src[beginOffsetLeft:endOffsetLeft])
						rightContent := string(src[beginOffsetRight:endOffsetRight])

						var sectionSegments [2][]*annotate.Annotation
						u7.HighlightedDiffFunc(leftContent, rightContent, &sectionSegments, [2]int{beginOffsetLeft, beginOffsetRight})

						anns = append(anns, sectionSegments[0]...)
						anns = append(anns, sectionSegments[1]...)
					}
					lastDel, lastIns = -1, -1
				}
			}

			sort.Sort(anns)

			out, err := annotate.Annotate(src, anns, func(w io.Writer, b []byte) { template.HTMLEscape(w, b) })
			if err != nil {
				return nil, false
			}

			return out, true
		}
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

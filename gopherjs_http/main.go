package gopherjs_http

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/go-on/gopherjslib"
)

// HtmlFile returns a handler that serves the given .html file, with the "text/go" scripts compiled to JavaScript via GopherJS.
func HtmlFile(name string) http.Handler {
	return &htmlFile{name: name}
}

type htmlFile struct {
	name string
}

func (this *htmlFile) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open(this.name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ProcessHtml(file).WriteTo(w)
}

// TODO: Write into writer, no need for buffer (unless want to be able to abort on error?). Or, alternatively, parse html and serve minified version?
func ProcessHtml(r io.Reader) *bytes.Buffer {
	insideTextGo := false
	tokenizer := html.NewTokenizer(r)
	var buf bytes.Buffer

	for {
		if tokenizer.Next() == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				return &buf
			}

			return &bytes.Buffer{}
		}

		token := tokenizer.Token()
		switch token.Type {
		case html.DoctypeToken:
			buf.WriteString(token.String())
		case html.CommentToken:
			buf.WriteString(token.String())
		case html.StartTagToken:
			if token.DataAtom == atom.Script && getType(token.Attr) == "text/go" {
				insideTextGo = true

				buf.WriteString(`<script type="text/javascript">`)

				if src := getSrc(token.Attr); src != "" {
					if data, err := ioutil.ReadFile(src); err == nil {
						buf.WriteString(goToJs(string(data)))
					}
				}
			} else {
				buf.WriteString(token.String())
			}
		case html.EndTagToken:
			if token.DataAtom == atom.Script && insideTextGo {
				insideTextGo = false
			}
			buf.WriteString(token.String())
		case html.SelfClosingTagToken:
			// TODO: Support <script type="text/go" src="..." />.
			buf.WriteString(token.String())
		case html.TextToken:
			if insideTextGo {
				buf.WriteString(goToJs(token.Data))
			} else {
				buf.WriteString(token.Data)
			}
		default:
			panic("unknown token type")
		}
	}
}

func getType(attrs []html.Attribute) string {
	for _, attr := range attrs {
		if attr.Key == "type" {
			return attr.Val
		}
	}
	return ""
}

func getSrc(attrs []html.Attribute) string {
	for _, attr := range attrs {
		if attr.Key == "src" {
			return attr.Val
		}
	}
	return ""
}

func goToJs(goCode string) (jsCode string) {
	started := time.Now()
	defer func() { fmt.Println("goToJs taken:", time.Since(started)) }()

	code := strings.NewReader(goCode)

	var out bytes.Buffer
	err := gopherjslib.Build(code, &out, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "console.error(\"" + template.JSEscapeString(err.Error()) + "\");"
	}

	return out.String()
}

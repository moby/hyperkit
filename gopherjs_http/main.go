package gopherjs_http

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/go-on/gopherjslib"
)

// StaticHtmlFile returns a handler that statically serves the given .html file, with the "text/go" script tags compiled to JavaScript via GopherJS.
//
// It reads file from disk and recompiles "text/go" script tags on startup only.
func StaticHtmlFile(name string) http.Handler {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return &staticHtmlFile{
		content: ProcessHtml(file).Bytes(),
	}
}

type staticHtmlFile struct {
	content []byte
}

func (this *staticHtmlFile) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(this.content)
}

// HtmlFile returns a handler that serves the given .html file, with the "text/go" script tags compiled to JavaScript via GopherJS.
//
// It reads file from disk and recompiles "text/go" script tags on every request.
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

// ProcessHtml takes HTML with "text/go" script tags and replaces them with compiled JavaScript script tags.
//
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

				if srcs := getSrcs(token.Attr); len(srcs) != 0 {
					buf.WriteString(handleJsError(goFilesToJs(srcs)))
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
				buf.WriteString(handleJsError(goToJs(token.Data)))
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

func getSrcs(attrs []html.Attribute) (srcs []string) {
	for _, attr := range attrs {
		if attr.Key == "src" {
			srcs = append(srcs, attr.Val)
		}
	}
	return srcs
}

func handleJsError(jsCode string, err error) string {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return `console.error("` + template.JSEscapeString(err.Error()) + `");`
	}
	return jsCode
}

func goFilesToJs(goFiles []string) (jsCode string, err error) {
	started := time.Now()
	defer func() { fmt.Println("goFilesToJs taken:", time.Since(started)) }()

	var out bytes.Buffer
	builder := gopherjslib.NewBuilder(&out, nil)

	for _, goFile := range goFiles {
		file, err := os.Open(goFile)
		if err != nil {
			return "", err
		}
		defer file.Close()

		builder.Add(goFile, file)
	}

	err = builder.Build()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func goToJs(goCode string) (jsCode string, err error) {
	started := time.Now()
	defer func() { fmt.Println("goToJs taken:", time.Since(started)) }()

	code := strings.NewReader(goCode)

	var out bytes.Buffer
	err = gopherjslib.Build(code, &out, nil)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

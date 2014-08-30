package gopherjs_http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/shurcooL/go/pipe_util"
	"gopkg.in/pipe.v2"
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
	processHtmlFile(file).WriteTo(w)
}

// TODO: Write into writer, no need for buffer. Or, alternatively, parse html and serve minified version?
// TODO: Support non-inline <script src="..."> tags.
func processHtmlFile(r io.Reader) *bytes.Buffer {
	var buff bytes.Buffer
	tokenizer := html.NewTokenizer(r)

	depth := 0

	for {
		if tokenizer.Next() == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				return &buff
			}

			return &bytes.Buffer{}
		}

		token := tokenizer.Token()
		switch token.Type {
		case html.DoctypeToken:
			buff.WriteString(token.String())
		case html.CommentToken:
			buff.WriteString(token.String())
		case html.StartTagToken:
			if token.DataAtom == atom.Script && getType(token.Attr) == "text/go" {
				depth++
				//goon.Dump(token.Attr)
				buff.WriteString(`<script type="text/javascript">`)
			} else {
				buff.WriteString(token.String())
			}
		case html.EndTagToken:
			if token.DataAtom == atom.Script && depth > 0 {
				depth--
			}
			buff.WriteString(token.String())
		case html.SelfClosingTagToken:
			buff.WriteString(token.String())
		case html.TextToken:
			if depth > 0 {
				buff.WriteString(goToJs(token.Data))
			} else {
				buff.WriteString(token.Data)
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

func goToJs(goCode string) (jsCode string) {
	started := time.Now()
	defer func() { fmt.Println("goToJs taken:", time.Since(started)) }()

	// TODO: Don't shell out, and avoid having to write/read temporary files, instead
	//       use http://godoc.org/github.com/gopherjs/gopherjs/compiler directly, etc.
	p := pipe.Script(
		pipe.Line(
			pipe.Print(goCode),
			pipe.WriteFile("tmp.go", 0666),
		),
		pipe.Exec("gopherjs", "build", "tmp.go"),
		pipe.ReadFile("tmp.js"),
	)

	// Use a temporary dir.
	tempDir, err := ioutil.TempDir("", "gopherjs_")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "warning: error removing temp dir:", err)
		}
	}()

	stdout, stderr, err := pipe_util.DividedOutputDir(p, tempDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Stderr.Write(stderr)
	}

	return string(stdout)
}

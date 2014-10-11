// Package markdown_http provides an http.Handler for serving Markdown over http.
package markdown_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shurcooL/go/u/u1"
)

// MarkdownHandlerFunc is an http.Handler that serves rendered Markdown.
type MarkdownHandlerFunc func(req *http.Request) (markdown []byte, err error)

func (this MarkdownHandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	markdown, err := this(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if _, plain := req.URL.Query()["plain"]; plain {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(markdown)
	} else if _, github := req.URL.Query()["github"]; github {
		w.Header().Set("Content-Type", "text/html")
		started := time.Now()
		u1.WriteGitHubFlavoredMarkdownViaGitHub(w, markdown)
		fmt.Println("rendered GFM via GitHub, took", time.Since(started))
	} else {
		w.Header().Set("Content-Type", "text/html")
		started := time.Now()
		u1.WriteGitHubFlavoredMarkdownViaLocal(w, markdown)
		fmt.Println("rendered GFM locally, took", time.Since(started))
	}
}

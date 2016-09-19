// Package openutil displays Markdown or HTML in a new browser tab.
package openutil

import (
	"net/http"
	"time"

	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	"github.com/shurcooL/go/gfmutil"
	"github.com/shurcooL/go/httpstoppable"
	"github.com/shurcooL/go/open"
)

// DisplayMarkdownInBrowser displays given Markdown in a new browser window/tab.
func DisplayMarkdownInBrowser(markdown []byte) {
	stopServerChan := make(chan struct{})

	handler := func(w http.ResponseWriter, req *http.Request) {
		gfmutil.WriteGitHubFlavoredMarkdownViaLocal(w, markdown)

		// TODO: A better way to fix: /assets/gfm/gfm.css Failed to load resource: net::ERR_CONNECTION_REFUSED.
		// HACK: Give some time for other assets to finish loading.
		go func() {
			time.Sleep(1 * time.Second)
			stopServerChan <- struct{}{}
		}()
	}

	http.HandleFunc("/index", handler)
	http.Handle("/assets/gfm/", http.StripPrefix("/assets/gfm", http.FileServer(gfmstyle.Assets))) // Serve the "/assets/gfm/gfm.css" file.
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	// TODO: Consider using httptest.NewServer.
	open.Open("http://localhost:7044/index")

	err := httpstoppable.ListenAndServe("localhost:7044", nil, stopServerChan)
	if err != nil {
		panic(err)
	}
}

// DisplayHTMLInBrowser displays given html page in a new browser window/tab.
// query can be empty, otherwise it should begin with "?" like "?key=value".
func DisplayHTMLInBrowser(mux *http.ServeMux, stopServerChan <-chan struct{}, query string) {
	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	open.Open("http://localhost:7044/index" + query)

	err := httpstoppable.ListenAndServe("localhost:7044", mux, stopServerChan)
	if err != nil {
		panic(err)
	}
}

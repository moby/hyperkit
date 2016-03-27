// Package openutil displays Markdown or HTML in a new browser tab.
package openutil

import (
	"net/http"

	"github.com/shurcooL/go/gfmutil"
	"github.com/shurcooL/go/httpstoppable"
	"github.com/shurcooL/go/open"
)

// DisplayMarkdownInBrowser displays given Markdown in a new browser window/tab.
func DisplayMarkdownInBrowser(markdown []byte) {
	stopServerChan := make(chan struct{})

	handler := func(w http.ResponseWriter, req *http.Request) {
		gfmutil.WriteGitHubFlavoredMarkdownViaLocal(w, markdown)

		stopServerChan <- struct{}{}
	}

	http.HandleFunc("/index", handler)
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
func DisplayHTMLInBrowser(mux *http.ServeMux, stopServerChan <-chan struct{}, query string) {
	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	open.Open("http://localhost:7044/index" + query)

	err := httpstoppable.ListenAndServe("localhost:7044", mux, stopServerChan)
	if err != nil {
		panic(err)
	}
}

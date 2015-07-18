// Package u3 displays Markdown or HTML in a new browser tab.
package u3

import (
	"net/http"

	"github.com/shurcooL/go/gists/gist7390843"
	"github.com/shurcooL/go/u/u1"
	"github.com/shurcooL/go/u/u4"
)

// Displays given Markdown in a new browser window/tab.
func DisplayMarkdownInBrowser(markdown []byte) {
	stopServerChan := make(chan struct{})

	handler := func(w http.ResponseWriter, req *http.Request) {
		u1.WriteGitHubFlavoredMarkdownViaLocal(w, markdown)

		stopServerChan <- struct{}{}
	}

	http.HandleFunc("/index", handler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	// TODO: Consider using httptest.NewServer.
	u4.Open("http://localhost:7044/index")

	err := gist7390843.ListenAndServeStoppable("localhost:7044", nil, stopServerChan)
	if err != nil {
		panic(err)
	}
}

// Displays given html page in a new browser window/tab.
func DisplayHtmlInBrowser(mux *http.ServeMux, stopServerChan <-chan struct{}, query string) {
	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	u4.Open("http://localhost:7044/index" + query)

	err := gist7390843.ListenAndServeStoppable("localhost:7044", mux, stopServerChan)
	if err != nil {
		panic(err)
	}
}

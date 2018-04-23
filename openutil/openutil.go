// Package openutil displays Markdown or HTML in a new browser tab.
package openutil

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	"github.com/shurcooL/go/gfmutil"
	"github.com/shurcooL/go/open"
)

// TODO: This code is extremely hacky and ridden with race conditions. :(
//       Ideally, it should be rewritten in a clean way, or deleted.
//       At this time, it's needed by goimporters and goimportgraph commands.

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

	// TODO: Start TCP listener before launching the browser to navigate to the page (else it's a race).
	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	// TODO: Consider using httptest.NewServer.
	open.Open("http://localhost:7044/index")

	server := &http.Server{Addr: "localhost:7044"}
	go func() {
		<-stopServerChan
		err := server.Close()
		if err != nil {
			log.Println("server.Close:", err)
		}
	}()
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(fmt.Errorf("server.ListenAndServe: %v", err))
	}
}

// DisplayHTMLInBrowser displays the /index HTML page of given mux in a new browser window/tab.
// query can be empty, otherwise it should begin with "?", like "?key=value".
func DisplayHTMLInBrowser(mux *http.ServeMux, stopServerChan <-chan struct{}, query string) {
	// TODO: Start TCP listener before launching the browser to navigate to the page (else it's a race).
	// TODO: Aquire a free port similarly to using ioutil.TempFile() for files.
	// TODO: Consider using httptest.NewServer.
	open.Open("http://localhost:7044/index" + query)

	server := &http.Server{Addr: "localhost:7044", Handler: mux}
	go func() {
		<-stopServerChan
		err := server.Close()
		if err != nil {
			log.Println("server.Close:", err)
		}
	}()
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(fmt.Errorf("server.ListenAndServe: %v", err))
	}
}

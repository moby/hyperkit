// Package raw_file_server provides a http.Handler that serves the given virtual file system without special handling of index.html.
package raw_file_server

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"path"
	"strings"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

type rawFileServer struct {
	// TODO: Use vfs.FileSystem.
	root http.FileSystem
}

// New returns a raw file server, that serves the given virtual file system without special handling of index.html.
func New(root vfs.FileSystem) http.Handler {
	// TODO: Use vfs.FileSystem.
	return &rawFileServer{httpfs.New(root)}
}

func (f *rawFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/") {
		r.URL.Path = "/" + r.URL.Path
	}
	serveFile(w, r, f.root, path.Clean(r.URL.Path))
}

func dirList(w http.ResponseWriter, f http.File, name string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", path.Clean(name+"/.."), "..")
	for {
		dirs, err := f.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()
			if d.IsDir() {
				name += "/"
			}
			// name may contain '?' or '#', which must be escaped to remain
			// part of the URL path, and not indicate the start of a query
			// string or fragment.
			url := url.URL{Path: name}
			fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), html.EscapeString(name))
		}
	}
	fmt.Fprintf(w, "</pre>\n")
}

// name is '/'-separated, not filepath.Separator.
func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string) {
	f, err := fs.Open(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to canonical path: / at end of directory url
	// r.URL.Path always begins with /
	url := r.URL.Path
	if d.IsDir() {
		if url[len(url)-1] != '/' {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}
	} else {
		if url[len(url)-1] == '/' {
			localRedirect(w, r, "../"+path.Base(url))
			return
		}
	}

	// A directory?
	if d.IsDir() {
		// TODO: Consider using checkLastModified?
		/*if checkLastModified(w, r, d.ModTime()) {
			return
		}*/
		dirList(w, f, name)
		return
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

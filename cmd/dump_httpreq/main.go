package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/shurcooL/go-goon"
)

var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address")

func dumpRequestHandler(w http.ResponseWriter, req *http.Request) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dump))
	goon.DumpExpr(req.URL.Query())
}

func main() {
	flag.Parse()

	fmt.Printf("Starting http request dumper, listening on %q...\n", *httpFlag)

	err := http.ListenAndServe(*httpFlag, http.HandlerFunc(dumpRequestHandler))
	if err != nil {
		panic(err)
	}
}

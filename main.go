package gist7390843

import (
	"net"
	"net/http"
	"strings"

	. "gist.github.com/5286084.git"
)

func ListenAndServeStoppable(addr string, handler http.Handler, stopServerChan <-chan bool) error {
	server := &http.Server{Addr: addr, Handler: handler}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	go func() { <-stopServerChan; err := listener.Close(); CheckError(err) }()
	err = server.Serve(listener)
	switch {
	case err == nil:
		panic("Supposed to get an error from Serve due to listener closed, but didn't...")
	case strings.Contains(err.Error(), "use of closed network connection"):
		return nil
	default:
		return err
	}
}

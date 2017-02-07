// Package httpstoppable provides ListenAndServe like http.ListenAndServe,
// but with ability to stop it.
//
// Deprecated: Go 1.8 adds native support for stopping a server in net/http.
// Once 1.8 is out, net/http should be used instead. This package will be
// removed shortly thereafter.
package httpstoppable

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// Handler is typically nil, in which case the http.DefaultServeMux is
// used.
//
// When receiving from stop unblocks (because it's closed or a value is sent),
// listener is closed and ListenAndServe returns with nil error.
// Otherise, it always returns a non-nil error.
func ListenAndServe(addr string, handler http.Handler, stop <-chan struct{}) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		<-stop
		err := ln.Close()
		if err != nil {
			log.Println("httpstoppable.ListenAndServe: error closing listener:", err)
		}
	}()
	err = srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	switch { // Serve always returns a non-nil error.
	case strings.Contains(err.Error(), "use of closed network connection"):
		return nil
	default:
		return err
	}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe so dead TCP connections
// (e.g. closing laptop mid-download) eventually go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

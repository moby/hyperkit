// Package httpstoppable offers ListenAndServe, like http.ListenAndServe, but with ability to stop it externally.
package httpstoppable

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// ListenAndServe is like http.ListenAndServe, but it closes listener socket when stop receives a value.
func ListenAndServe(addr string, handler http.Handler, stop <-chan struct{}) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	ln, err := net.Listen("tcp", srv.Addr)
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
	switch {
	case err == nil:
		panic("Supposed to get an error from Serve due to listener closed, but didn't...")
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

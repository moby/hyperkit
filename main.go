package gist7390843

import (
	"net"
	"net/http"
	"strings"
)

// Closes listener socket when stopServerChan receives a value.
func ListenAndServeStoppable(addr string, handler http.Handler, stopServerChan <-chan struct{}) error {
	server := &http.Server{Addr: addr, Handler: handler}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	go func() {
		<-stopServerChan
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}()
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

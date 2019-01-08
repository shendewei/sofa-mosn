package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/alipay/sofa-mosn/pkg/module/http2"
)

func main() {
	h2Server := &http2.Server{IdleTimeout: 1 * time.Minute}
	handler := &HTTPHandler{}
	s := Server{
		Server:  h2Server,
		Handler: handler,
	}
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return
		}
		go s.Serve(conn)
	}
}

type Server struct {
	Server  *http2.Server
	Handler http.Handler
}

func (s *Server) Serve(conn net.Conn) {
	s.Server.ServeConn(conn, &http2.ServeConnOpts{Handler: s.Handler})
}

type HTTPHandler struct{}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[UPSTREAM]receive request %s\n", r.URL)

	// read body
	buf := make([]byte, 1024)
	r.Body.Read(buf)
	fmt.Println("Receive Data: ", string(buf))

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "RequestURI: %q\n", r.RequestURI)
	fmt.Fprintf(w, "URL: %#v\n", r.URL)
	fmt.Fprintf(w, "Body.ContentLength: %d (-1 means unknown)\n", r.ContentLength)
	fmt.Fprintf(w, "Close: %v (relevant for HTTP/1 only)\n", r.Close)
	fmt.Fprintf(w, "TLS: %#v\n", r.TLS)
	fmt.Fprintf(w, "\nHeaders:\n")

	r.Header.Write(w)
}

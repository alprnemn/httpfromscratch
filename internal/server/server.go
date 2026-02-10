package server

import (
	"fmt"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"log"
	"net"
	"sync/atomic"
)

// Server represents a minimal HTTP server implementation built
// directly on top of raw TCP sockets. It manages the listening socket,
// connection lifecycle, and request dispatching.
type Server struct {
	Addr     string
	listener net.Listener
	closed   atomic.Bool
	Handler  Handler
}

// Serve starts a new TCP listener on the given port and begins accepting
// incoming connections. For each accepted connection, the provided
// handler is invoked in its own goroutine.
//
// It returns a Server instance that can be used to manage the server
// lifecycle, including graceful shutdown.
func Serve(port uint16, handler Handler) (*Server, error) {

	portStr := fmt.Sprintf(":%d", port)

	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Addr:     portStr,
		Handler:  handler,
		listener: l,
	}

	go s.listen()

	return s, nil
}

// Close gracefully shuts down the server by marking it as closed
// and closing the underlying TCP listener.
func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

// listen continuously accepts incoming TCP connections and dispatches
// each connection to a separate goroutine for request handling.
// The loop exits when the server is marked as closed or the listener
// encounters a fatal error.
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if s.closed.Load() {
			return
		}
		if err != nil {
			return
		}

		fmt.Println("connection accepted: ", conn.RemoteAddr())

		go s.handle(conn)

	}

}

// handle processes a single client connection. It parses the incoming
// HTTP request, initializes a response writer, and invokes the
// configured handler to generate the response.
//
// The connection is closed automatically when request handling completes.
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.ParseRequest(conn)
	if err != nil {
		log.Fatal(err)
	}

	writer := response.NewWriter(conn)
	s.Handler(writer, req)

}

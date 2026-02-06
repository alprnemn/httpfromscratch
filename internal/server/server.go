package server

import (
	"fmt"
	"httpfromscratch/internal/request"
	r "httpfromscratch/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type ConnState int

const (
	StateNew ConnState = iota
	StateActive
	StateClosed
)

type Server struct {
	Addr     string
	listener net.Listener
	closed   atomic.Bool
	Handler  Handler
}

// Serve func serves http
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

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

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

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.ParseRequest(conn)
	if err != nil {
		log.Fatal(err)
	}

	writer := r.NewWriter(conn)
	s.Handler(writer, req)

}

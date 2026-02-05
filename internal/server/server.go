package server

import (
	"fmt"
	r "httpfromscratch/internal/response"
	"io"
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
}

// Serve func serves http
func Serve(port uint16) (*Server, error) {

	s := &Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return nil, err
	}

	s.listener = l

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

func (s *Server) handle(conn io.ReadWriteCloser) {

	if err := r.WriteStatusLine(conn, 400); err != nil {
		return
	}
	conn.Close()
}

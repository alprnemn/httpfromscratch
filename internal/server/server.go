package server

import (
	"fmt"
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
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!")

	conn.Write(out)
	fmt.Println("conn writed")
	conn.Close()
}

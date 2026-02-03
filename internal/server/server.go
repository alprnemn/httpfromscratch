package server

import (
	"fmt"
	rq "httpfromscratch/internal/request"
	"log"
	"net"
	"strconv"
)

type Server struct {
	Addr      string
	ConnState ConnState
	Listener  net.Listener
}

type ConnState int

const (
	StateNew ConnState = iota
	StateActive
	StateClosed
)

func Serve(port int) (*Server, error) {

	sv := &Server{
		Addr:      ":" + strconv.Itoa(port),
		ConnState: StateNew,
	}

	l, err := net.Listen("tcp", sv.Addr)
	if err != nil {
		log.Fatal(err)
	}

	sv.Listener = l

	go sv.listen()

	return sv, nil
}

func (s *Server) Close() error {
	s.ConnState = StateClosed
	s.Listener.Close()
	return nil
}

func (s *Server) listen() {

	s.ConnState = StateActive

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			continue
		}

		fmt.Println("connection accepted", conn.RemoteAddr())

		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {
	req, err := rq.ParseRequest(conn)
	if err != nil {
		log.Fatal("error: ", err)
	}

	rq.WriteReqLines(req)

	req.Headers.ForEach(func(n, v string) {
		fmt.Println(n, v)
	})

	fmt.Println("connection closed", conn.RemoteAddr())

}

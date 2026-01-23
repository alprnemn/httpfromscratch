package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {

	//l, err := net.Listen("tcp", "localhost:41234")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for {
	//	conn, err := l.Accept()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Println("connection accepted for: %s", conn.RemoteAddr())
	//	go handleConnection(conn)
	//}

	data := []byte("\r\n")
	n := bytes.Index(data, []byte("\r\n"))
	fmt.Println("n: ", n)

}

func handleConnection(conn net.Conn) {
	req, err := ParseRequest(conn)
	if err != nil {
		log.Println("err: ", err)
	}
	fmt.Println(req)
}

// GET / HTTP/1.1\r\nHost:localhost:3000\r\nAgent:Curl2.3\r\n\r\nbodystarthere

func ParseRequest(conn net.Conn) (*Request, error) {

	request := NewRequest()

	buf := make([]byte, BufferSize)

	bufIndex := 0

	for !request.done() {
		n, err := conn.Read(buf[bufIndex:])
		if err != nil {
			return nil, err
		}

		bufIndex += n

		read, err := request.Parse(buf[:bufIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[read:bufIndex])

		bufIndex -= read
	}
	return request, nil
}

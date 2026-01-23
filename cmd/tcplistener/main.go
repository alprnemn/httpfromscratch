package main

import (
	"fmt"
	rq "httpfromscratch/internal/request"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		fmt.Println("connection accepted")

		go handleConnection(conn)

	}

}

// handle connection and get lines from the connection
func handleConnection(conn net.Conn) {

	req, err := rq.ParseRequest(conn)
	if err != nil {
		log.Fatal("error: ", err)
	}

	fmt.Println("ip adress from conn: ", conn.RemoteAddr())
	rq.WriteReqLines(req)

	req.Headers.ForEach(func(n, v string) {
		fmt.Println(n, v)
	})

	fmt.Println("connection closed")
}

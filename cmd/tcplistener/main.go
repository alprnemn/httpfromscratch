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
	rq.WriteHeaders(req)
	rq.WriteReqLines(req)

	fmt.Println("connection closed")
}

//lines := getLinesChannel(conn)
//for line := range lines {
//	fmt.Println(line)
//}
// // this function reads lines from the string channel
//func getLinesChannel(conn net.Conn) <-chan string {
//
//	ch := make(chan string)
//
//	go func() {
//		defer close(ch)
//		defer conn.Close()
//
//		currentLineContents := ""
//
//		for {
//			// create buffer
//			buf := make([]byte, 8)
//
//			// read buffer from connection
//			n, err := conn.Read(buf)
//
//			if err != nil {
//				// if raw is ended throws eof and catch it to add last raw
//				if errors.Is(err, io.EOF) {
//					if currentLineContents != "" {
//						ch <- currentLineContents
//					}
//					return
//				}
//				fmt.Println(err)
//				return
//			}
//
//			// split into parts
//			parts := strings.Split(string(buf[:n]), "\n")
//
//			// send parts to the channel except last one
//			for i := 0; i < len(parts)-1; i++ {
//				ch <- currentLineContents + parts[i]
//				currentLineContents = ""
//			}
//			currentLineContents += parts[len(parts)-1]
//		}
//	}()
//
//	return ch
//}

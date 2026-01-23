package main

import (
	"fmt"
	r "httpfromscratch/internal/request"
	"io"
)

func main() {

	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 5,
	}

	request, err := r.ParseRequest(reader)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	r.WriteReqLines(request)
	fmt.Println(request.Headers.Get("Host"))
	fmt.Println(request.Headers.Get("User-Agent"))
	fmt.Println(request.Headers.Get("Accept"))

}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// [4, 4, 3 , 6]
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}

	endIndex := cr.pos + cr.numBytesPerRead

	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}

	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil

}

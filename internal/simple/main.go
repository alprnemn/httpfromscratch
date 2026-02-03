package main

import (
	"fmt"
	r "httpfromscratch/internal/request"
	"io"
)

func main() {

	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:4069\r\n" + // 25
			"Content-Length:111\r\n" + // 20
			"\r\n" +
			"{ name : 'alperen'}",

		numBytesPerRead: 3,
	}
	request, err := r.ParseRequest(reader)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	r.WriteReqLines(request)
	request.Headers.ForEach(func(n, v string) {
		fmt.Println(n, v)
	})

	fmt.Println("req body: ", string(request.Body))

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

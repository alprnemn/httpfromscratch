package main

import (
	"errors"
	"fmt"
	cmn "httpfromscratch/common"
	"io"
	"log"
	"strings"
)

var validMethods = map[string]struct{}{
	"GET":    {},
	"POST":   {},
	"PUT":    {},
	"DELETE": {},
	"PATCH":  {},
}

type parserState string

const BufferSize = 8

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

func main() {

	reader := &chunkReader{
		data:            "GET /asd HTTP/1.1\r\nHTPPHEADERASD:ASFASF\r\n",
		numBytesPerRead: 3,
	}

	req, err := ParseRequestFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	WriteReqLines(req)

}

// asdasda\r\nHTPPHEADERASD:ASFASF\r\n

// 3 asd
// 6 asd
// 9 a\r
// 12 \nH buf = asdasda\r\nH

// ParseRequestFromReader parses entire http request
func ParseRequestFromReader(reader io.Reader) (*Request, error) {

	request := NewRequest()

	buf := make([]byte, BufferSize)

	bufIndex := 0

	for !request.done() {

		if bufIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// 3 tane okur ve buf a yazar
		n, err := reader.Read(buf[bufIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = StateDone
				break
			}
			return nil, err
		}

		bufIndex += n

		readN, err := request.parse(buf[:bufIndex])
		if err != nil {
			return nil, err
		}
		fmt.Println("[readn:]", buf[readN:bufIndex])
		copy(buf, buf[readN:bufIndex])
		bufIndex -= readN

	}

	fmt.Println("last buf: ", string(buf))
	return request, nil
}

// parseRequestLine parses request line parts and returns *RequestLine readed n and err
func parseRequestLine(b []byte) (*RequestLine, int, error) {

	crlfIndex := strings.Index(string(b), "\r\n")

	if crlfIndex == -1 {
		return nil, 0, nil
	}

	line := string(b[:crlfIndex])
	read := crlfIndex + len([]byte("\r\n"))
	parts := strings.Fields(line)

	if err := validateParts(parts); err != nil {
		return nil, 0, err
	}

	reqline := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][5:],
	}

	return reqline, read, nil
}

// validateParts validates RequestLine Parts
func validateParts(parts []string) error {
	if len(parts) != 3 {
		return errors.New("request line parts must be 3")
	}

	if err := validateMethod(parts[0]); err != nil {
		return err
	}

	if err := validateTarget(parts[1]); err != nil {
		return err
	}

	if err := validateHTTPVersion(parts[2]); err != nil {
		return err
	}

	return nil
}

// validateMethod check method is valid
func validateMethod(m string) error {
	if m == "" {
		return cmn.ErrInvalidMethodFormat
	}

	if !cmn.IsAllUpper(m) {
		return errors.New("method chars must be capitalize")
	}

	if _, ok := validMethods[m]; !ok {
		return cmn.ErrMethodNotImplemented
	}

	return nil
}

// validateTarget check request target is valid
func validateTarget(t string) error {
	if t == "" {
		return errors.New("target cannot be empty")
	}

	if t[0] != '/' {
		return errors.New("target must start /")
	}

	return nil
}

// validateHTTPVersion valids version
func validateHTTPVersion(s string) error {
	if s == "" {
		return errors.New("version cannot be empty")
	}

	versionParts := strings.Split(s, "/")

	if len(versionParts) != 2 {
		return errors.New("invalid http version format")
	}

	if versionParts[0] != "HTTP" {
		return errors.New("invalid http version format")
	}

	if versionParts[1] != "1.1" {
		return errors.New("only supported version is 1.1")
	}
	return nil
}

// Writes request line fields to console
func WriteReqLines(req *Request) {
	fmt.Println("req line method: ", req.RequestLine.Method)
	fmt.Println("req line target: ", req.RequestLine.RequestTarget)
	fmt.Println("req line http version: ", req.RequestLine.HttpVersion)
}

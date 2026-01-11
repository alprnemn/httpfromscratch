package main

import (
	"errors"
	"fmt"
	"httpfromscratch/common"
	"io"
	"strings"
)

type Request struct {
	RequestLine *RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func main() {
	r := strings.NewReader("GET /products HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")

	req, err := RequestFromReader(r)
	common.FailOnError(err, "errorrrrr")

	fmt.Println("req method: ", req.RequestLine.Method)
	fmt.Println("req target: ", req.RequestLine.RequestTarget)
	fmt.Println("http version: ", req.RequestLine.HttpVersion)
}

func RequestFromReader(r io.Reader) (*Request, error) {

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(b), "\r\n")

	req := &Request{}

	reqLine, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, err
	}

	if req.RequestLine != nil {
		req.RequestLine = reqLine

	}

	return req, nil
}
func parseRequestLine(s string) (*RequestLine, error) {

	if s == "" || len(s) == 0 {
		return nil, errors.New("b cannot be empty")
	}

	parts := strings.Split(s, " ")

	reqline := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][5:],
	}

	return reqline, nil
}

//
//func RequestFromReader(r io.Reader) (*Request, error) {
//
//	b, err := io.ReadAll(r)
//	common.FailOnError(err, "error io read all from r")
//	req := &Request{}
//
//	parts := strings.Split(string(b), "\r\n")
//
//	for i := 0; i < len(b); i++ {
//		if b[i] == '\r' && b[i+1] == '\n' {
//			reqLine, err := parseRequestLine(b[:i])
//			common.FailOnError(err, "error while parsing request line")
//			req.RequestLine = reqLine
//			break
//		}
//	}
//
//	return req, nil
//}

//func parseRequestLine(b []byte) (*RequestLine, error) {
//
//	if len(b) == 0 {
//		return nil, errors.New("b cannot be empty")
//	}
//
//	reqline := &RequestLine{}
//
//	state := "METHOD"
//
//	secIndex := 0
//
//	for i := 0; i < len(b); i++ {
//		switch state {
//		case "METHOD":
//			if b[i] == 32 {
//				reqline.Method = string(b[:i])
//				state = "TARGET"
//				secIndex = i + 1
//			}
//		case "TARGET":
//			if b[i] == 32 {
//				reqline.RequestTarget = string(b[secIndex:i])
//				reqline.HttpVersion = string(b[i+1:][5:])
//				break
//			}
//		}
//	}
//	return reqline, nil
//}

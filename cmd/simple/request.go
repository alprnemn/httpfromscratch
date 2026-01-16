package main

import (
	"bytes"
	"errors"
	"fmt"
	cmn "httpfromscratch/common"
	"io"
)

var validMethods = map[string]struct{}{
	"GET":    {},
	"POST":   {},
	"PUT":    {},
	"DELETE": {},
}

type parserState string

const BufferSize = 1024

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	State       parserState
}

func NewRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0

outer:
	for {
		switch r.State {
		case StateError:
			return 0, errors.New("error in request state")
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])

			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.State = StateDone

		case StateDone:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

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

		// reads 3(n) and writes to buf
		n, err := reader.Read(buf[bufIndex:])

		bufIndex += n

		readN, err := request.parse(buf[:bufIndex])
		if err != nil {
			return nil, err
		}

		fmt.Println("[readn:]", string(buf[readN:bufIndex]))
		fmt.Println("buf: ", string(buf))
		copy(buf, buf[readN:bufIndex])
		fmt.Println("buf2: ", string(buf))
		bufIndex -= readN
	}
	return request, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	crlf := []byte("\r\n")

	crlfIndex := bytes.Index(b, crlf)

	if crlfIndex == -1 {
		return nil, 0, nil
	}

	line := b[:crlfIndex]
	read := crlfIndex + len(crlf)

	parts := bytes.Fields(line)

	if err := validateParts(parts); err != nil {
		return nil, 0, err
	}

	reqline := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(parts[2][5:]),
	}

	return reqline, read, nil
}

func validateParts(parts [][]byte) error {
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

func validateMethod(m []byte) error {
	if len(m) == 0 {
		return cmn.ErrInvalidMethodFormat
	}

	if !cmn.IsAllUpperBytes(m) {
		return errors.New("method chars must be capitalize")
	}

	if _, ok := validMethods[string(m)]; !ok {
		return cmn.ErrMethodNotImplemented
	}

	return nil
}

func validateTarget(t []byte) error {
	if len(t) == 0 {
		return errors.New("target cannot be empty")
	}

	if t[0] != '/' {
		return errors.New("target must start /")
	}

	return nil
}

func validateHTTPVersion(v []byte) error {
	if len(v) == 0 {
		return errors.New("version cannot be empty")
	}

	parts := bytes.Split(v, []byte("/"))

	if len(parts) != 2 {
		return errors.New("invalid http version format")
	}

	if !bytes.Equal(parts[0], []byte("HTTP")) {
		return errors.New("invalid http version format")
	}

	if !bytes.Equal(parts[1], []byte("1.1")) {
		return errors.New("only supported version is 1.1")
	}

	return nil
}

func WriteReqLines(req *Request) {
	fmt.Println("req line method: ", req.RequestLine.Method)
	fmt.Println("req line target: ", req.RequestLine.RequestTarget)
	fmt.Println("req line http version: ", req.RequestLine.HttpVersion)
}

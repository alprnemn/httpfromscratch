package request

import (
	"bytes"
	"errors"
	"fmt"
	cmn "httpfromscratch/common"
	"httpfromscratch/internal/headers"
	"io"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        []byte
	State       parserState
}

func NewRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// ParseRequest parses entire http request
func ParseRequest(reader io.Reader) (*Request, error) {

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

		copy(buf, buf[readN:bufIndex])

		bufIndex -= readN
	}

	return request, nil

}

func (r *Request) parse(data []byte) (int, error) {

	read := 0

outer:
	for {

		currentData := data[read:]

		switch r.State {
		case StateError:
			return 0, errors.New("error in request state")
		case StateInit:
			rl, n, err := parseRequestLine(currentData)

			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.State = StateHeaders
			fmt.Println(rl.Method)
			fmt.Println(rl.RequestTarget)
			fmt.Println(rl.HttpVersion)
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				r.State = StateDone
			}

		case StateDone:
			break outer
		default:
			panic("something went wrongg at the req parse")
		}

	}

	return read, nil
}

// parseRequestLine parses request line parts and returns *RequestLine readed n and err
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

// validateParts validates RequestLine Parts
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

// validateMethod check method is valid
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

// validateTarget check request target is valid
func validateTarget(t []byte) error {
	if len(t) == 0 {
		return errors.New("target cannot be empty")
	}

	if t[0] != '/' {
		return errors.New("target must start /")
	}

	return nil
}

// validateHTTPVersion v0alids version
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

// Writes request line fields to consol
func WriteReqLines(req *Request) {
	fmt.Println("req line method: ", req.RequestLine.Method)
	fmt.Println("req line target: ", req.RequestLine.RequestTarget)
	fmt.Println("req line http version: ", req.RequestLine.HttpVersion)
}

//
//func WriteHeaders(req *Request) {
//
//	for k, v := range req.Headers. {
//		fmt.Printf("key:%s, value:%s\n", k, v)
//	}
//
//}

package request

import (
	"bytes"
	"errors"
	"fmt"
	cmn "httpfromscratch/common"
	"httpfromscratch/internal/headers"
	"io"
	"strconv"
)

// Request represents a parsed HTTP/1.1 request. It contains the request
// line, headers, body, and the internal parser state.
type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        string
	State       parserState
}

// NewRequest initializes Request struct using StateInit headers
// and empty body for Request struct
func NewRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
	}
}

// done reports whether the request parsing process has completed
// either successfully or due to an error.
func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

// hasBody determines whether the request contains a body by inspecting
// the Content-Length header.
func (r *Request) hasBody() bool {
	length := getInt(r.Headers, "content-length", 0)
	return length > 0
}

// RequestLine represents the first line of an HTTP request,
// consisting of method, request target, and HTTP version.
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// ParseRequest reads from the provided io.Reader and parses an entire
// HTTP/1.1 request. It blocks until the full request is received or
// an error occurs.
func ParseRequest(reader io.Reader) (*Request, error) {

	request := NewRequest()

	buf := make([]byte, BufferSize)

	bufIndex := 0

	totalBytesParsed := 0

	for !request.done() {
		if bufIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// reads (n) and writes to buf
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

		totalBytesParsed += readN

		copy(buf, buf[readN:bufIndex])

		bufIndex -= readN
	}

	return request, nil

}

// parse consumes raw byte data and advances the internal parser state
// machine until no more progress can be made or the request is complete.
// It returns the number of bytes successfully consumed.
func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
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
				if r.hasBody() {
					r.State = StateBody
				} else {
					r.State = StateDone
				}
			}

		case StateBody:
			length := getInt(r.Headers, "content-length", 0)

			if length == 0 {
				panic("chunk not implemented")
			}

			remaining := min(length-len(r.Body), len(currentData))

			r.Body += string(currentData[:remaining])

			read += remaining

			if len(r.Body) == length {
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

// parseRequestLine parses the HTTP request line from the provided byte slice.
// It returns the parsed RequestLine, the number of bytes consumed,
// and any parsing error.
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

// validateParts validates the three components of an HTTP request line:
// method, request target, and HTTP version.
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

// validateMethod validates the HTTP method token and ensures that
// it is uppercase and implemented by the server.
func validateMethod(m []byte) error {
	if len(m) == 0 {
		return cmn.ErrInvalidMethodFormat
	}

	if !cmn.IsAllUpperBytes(m) {
		return errors.New("method chars must be capitalize")
	}

	if _, ok := validMethods[string(m)]; !ok {
		return errors.New("method name not implemented")
	}

	return nil
}

// validateTarget validates the request target and ensures that it
// begins with a forward slash (/).
func validateTarget(t []byte) error {
	if len(t) == 0 {
		return errors.New("target cannot be empty")
	}

	if t[0] != '/' {
		return errors.New("target must start /")
	}

	return nil
}

// validateHTTPVersion validates the HTTP version token and ensures
// that only HTTP/1.1 is accepted.
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

// WriteReqLines prints the parsed request line fields to stdout.
// Intended primarily for debugging purposes.
func WriteReqLines(req *Request) {
	fmt.Println("req line method: ", req.RequestLine.Method)
	fmt.Println("req line target: ", req.RequestLine.RequestTarget)
	fmt.Println("req line http version: ", req.RequestLine.HttpVersion)
}

// getInt retrieves an integer header value from the headers collection.
// If the header does not exist or cannot be parsed, defaultValue is returned.
func getInt(headers *headers.Headers, name string, defaultValue int) int {

	valueStr, exists := headers.Header(name)
	if !exists {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return valueInt

}

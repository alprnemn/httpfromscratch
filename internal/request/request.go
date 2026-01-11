package request

import (
	"errors"
	cmn "httpfromscratch/common"
	"io"
	"strings"
)

var validMethods = map[string]struct{}{
	"GET":    {},
	"POST":   {},
	"PUT":    {},
	"DELETE": {},
	"PATCH":  {},
}

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// ParseRequestFromReader parses entire http request
func ParseRequestFromReader(r io.Reader) (*Request, error) {

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(b), "\r\n")

	reqLine, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, err
	}

	req := &Request{RequestLine: *reqLine}

	return req, nil
}

// parseRequestLine parses request line parts and returns *RequestLine
func parseRequestLine(s string) (*RequestLine, error) {

	if s == "" || len(s) == 0 {
		return nil, errors.New("s cannot be empty")
	}

	parts := strings.Fields(s)

	if err := validateParts(parts); err != nil {
		return nil, err
	}

	reqline := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][5:],
	}

	return reqline, nil
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

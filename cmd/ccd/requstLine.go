package main

import (
	"bytes"
	"errors"
	cmn "httpfromscratch/common"
)

type RequestLine struct {
	Method  string
	Target  string
	Version string
}

func NewRequsetLine() *RequestLine {
	return &RequestLine{}
}

func ParseRequestLine(data []byte) (*RequestLine, int, error) {

	crlf := []byte("\r\n")

	crlfIndex := bytes.Index(data, crlf)

	if crlfIndex == -1 {
		return nil, 0, nil
	}

	line := data[:crlfIndex]
	read := crlfIndex + len(crlf)

	parts := bytes.Fields(line)

	if err := validateParts(parts); err != nil {
		return nil, 0, err
	}

	rl := &RequestLine{
		Method:  string(parts[0]),
		Target:  string(parts[1]),
		Version: string(parts[2][5:]),
	}

	return rl, read, nil
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

var validMethods = map[string]struct{}{
	"GET":    {},
	"POST":   {},
	"PUT":    {},
	"DELETE": {},
	"PATCH":  {},
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

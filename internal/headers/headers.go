package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) (string, bool) {
	val, ok := h.headers[strings.ToLower(name)]
	return val, ok
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if _, ok := h.headers[name]; ok {
		h.headers[name] += "," + value
		return
	}
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

func (h *Headers) Parse(data []byte) (int, bool, error) {

	crlf := []byte("\r\n")
	read := 0
	done := false

	for {

		crLFidx := bytes.Index(data[read:], crlf)
		if crLFidx == -1 {
			break
		}

		if crLFidx == 0 {
			done = true
			read += len(crlf)
			break
		}

		name, value, err := parseHeader(data[read : read+crLFidx])

		if err != nil {
			return 0, false, err
		}
		read += crLFidx + len(crlf)

		h.Set(name, value)
	}

	return read, done, nil

}

// TODO: Validate headers ' ' by RFC make whitespaces appropriate
func parseHeader(header []byte) (string, string, error) {

	header = bytes.TrimSpace(header)

	parts := bytes.SplitN(header, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid header format (missing colon)")
	}

	if len(parts[1]) == 0 {
		return "", "", errors.New("header got no value")
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])

	if !isValidHeaderKey(key) {
		return "", "", fmt.Errorf("invalid header key: %s", string(key))
	}

	if !isValidHeaderValue(value) {
		return "", "", fmt.Errorf("invalid header value: %s", string(value))
	}

	return string(key), string(value), nil
}

// isValidHeaderValue checks key is valid
func isValidHeaderValue(value []byte) bool {
	for _, b := range value {
		// CR / LF
		if b == '\r' || b == '\n' {
			return false
		}

		// Control chars (except TAB)
		if b < 0x20 && b != '\t' {
			return false
		}

		// DEL
		if b == 0x7F {
			return false
		}
	}
	return true
}

// isValidHeaderKey checks value is valid
func isValidHeaderKey(key []byte) bool {
	if len(key) == 0 {
		return false
	}

	for _, b := range key {
		if !isTChar(b) {
			return false
		}
	}

	return true
}

// isTChar checks byte is valid (A-Z,a-z,0-9,validchars)
func isTChar(b byte) bool {
	// A-Z
	if b >= 'A' && b <= 'Z' {
		return true
	}

	// a-z
	if b >= 'a' && b <= 'z' {
		return true
	}

	// 0-9
	if b >= '0' && b <= '9' {
		return true
	}

	switch b {
	case '!', '#', '$', '%', '&', '\'', '*',
		'+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}

	return false
}

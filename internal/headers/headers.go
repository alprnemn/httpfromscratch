package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Headers represents a collection of HTTP header fields.
// Header names are stored in lowercase to allow for case-insensitive lookup.
type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

// Header Get retrieves the value of the specified header field.
// The lookup is case-insensitive. The boolean return value indicates
// whether the header exists.
func (h *Headers) Header(name string) (string, bool) {
	val, ok := h.headers[strings.ToLower(name)]
	return val, ok
}

// SetHeader Set inserts or updates the specified header field.
// Header names are normalized to lowercase for consistent storage.
//
// If the header already exists (except for Content-Length), the new value
// is appended using a comma, following HTTP header field semantics.
func (h *Headers) SetHeader(name, value string) {
	name = strings.ToLower(name)

	if name == "content-length" {
		h.headers[name] = value
		return
	}

	if _, ok := h.headers[name]; ok {
		h.headers[name] += "," + value
		return
	}
	h.headers[strings.ToLower(name)] = value
}

// ForEach iterates over all stored headers and invokes the provided
// callback function for each name-value pair.
func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

// Parse incrementally parses HTTP header fields from the provided byte slice.
// It returns the number of bytes consumed, a boolean indicating whether
// the header section is complete, and any parsing error encountered.
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

		h.SetHeader(name, value)
	}

	return read, done, nil

}

// parseHeader parses a single HTTP header line in the form:
//
//	<field-name>: <field-value>
//
// It validates both the header name and value according to RFC 9110 - Tokens 5.6.2
// and returns the parsed name and value.
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

// isValidHeaderValue validates an HTTP header field value according
// to RFC 9110. It rejects control characters, CR, LF, and DEL bytes.
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

// isValidHeaderKey validates an HTTP header field name according
// to RFC 9110 token rules.
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

package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) ParseHeader(data []byte) (n int, done bool, err error) {

	crLf := []byte("\r\n")

	endOfHeaders := []byte("\r\n\r\n")

	endCRLFIndex := bytes.Index(data, endOfHeaders)
	if endCRLFIndex == -1 {
		return 0, false, nil
	}

	allHeaders := data[:endCRLFIndex]
	read := endCRLFIndex + len(endOfHeaders)

	headers := bytes.Split(allHeaders, crLf)

	for i := 0; i < len(headers); i++ {
		k, v, err := validateHeader(headers[i])
		if err != nil {
			return 0, false, err
		}
		h[k] = v
	}

	return read, true, nil
}

// TODO: Validate headers ' ' by RFC make whitespaces appropriate
func validateHeader(header []byte) (string, string, error) {

	header = bytes.TrimSpace(header)

	parts := bytes.SplitN(header, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid header format (missing colon)")
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])

	if !isValidHeaderKey(key) {
		return "", "", fmt.Errorf("invalid header key: %s", string(key))
	}

	if !isValidHeaderValue(value) {
		return "", "", fmt.Errorf("invalid header value: %s", string(value))
	}

	key = bytes.ToLower(key)

	return string(key), string(value), nil
}

// isValidHeaderValue checks key is valid
func isValidHeaderValue(value []byte) bool {
	for _, b := range value {
		// CR / LF
		if b == '\r' || b == '\n' {
			return false
		}

		// Control chars
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

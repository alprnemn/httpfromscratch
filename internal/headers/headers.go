package headers

import (
	"bytes"
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

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if _, ok := h.headers[name]; ok {
		h.headers[name] += "," + value
		return
	}
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	// Host:asdf\r\nAgent
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

		name, value, err := parseHeader(data[read:crLFidx])
		if err != nil {
			return 0, false, err
		}

		read = crLFidx + len(crlf)

		h.Set(name, value)

	}

	return read, done, nil

}

func parseHeadersTrimeagen(fieldLine []byte) (string, string, error) {

	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil

}

// TODO: Validate headers ' ' by RFC make whitespaces appropriate
func parseHeader(header []byte) (string, string, error) {

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

	return string(key), string(value), nil
}

// isValidHeaderValue checks key is valid
func isValidHeaderValue(value []byte) bool {
	for _, b := range value {
		// CR / LF
		if b == '\r' || b == '\n' {
			return false
		}

		// Space
		if b == ' ' { // or b == 0x20
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

//func MyParseHeadersMethod(headers []byte) {
//	endOfHeaders := []byte("\r\n\r\n")
//
//	endCRLFIndex := bytes.Index(data, endOfHeaders)
//	if endCRLFIndex == -1 {
//		return 0, false, nil
//	}
//
//	allHeaders := data[:endCRLFIndex]
//	read := endCRLFIndex + len(endOfHeaders)
//
//	headers := bytes.Split(allHeaders, crlf)
//
//	for i := 0; i < len(headers); i++ {
//		k, v, err := validateHeader(headers[i])
//		if err != nil {
//			return 0, false, err
//		}
//		if _, ok := h[k]; ok {
//			h[k] = h[k] + "," + v
//		} else {
//			h[k] = v
//		}
//	}
//}

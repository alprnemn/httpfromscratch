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

	//if err := validateHeaders(headers); err != nil {
	//	return 0, false, err
	//}

	for i := 0; i < len(headers); i++ {
		l, r, err := validateHeader(headers[i])
		if err != nil {
			return 0, false, err
		}
		h[l] = r
	}

	return read, true, nil
}

// TODO: Validate headers ' ' by RFC make whitespaces appropriate
func validateHeader(header []byte) (string, string, error) {

	header = bytes.TrimSpace(header)

	parts := bytes.SplitN(header, []byte(":"), 2)

	if bytes.Index(parts[0], []byte(" ")) != -1 {
		return "", "", fmt.Errorf("key cannot take ' ' space at %s field-name", string(parts[0]))
	}

	parts[1] = bytes.TrimSpace(parts[1])

	left := string(parts[0])
	right := string(parts[1])

	return left, right, nil
}

//POST /coffee HTTP/1.1\r\n
//    Host: localhost:42069     \r\n
//User-Agent: curl/7.81.0\r\n
//Accept: */*\r\n
//Content-Length: 11\r\n
//\r\n
//hello world

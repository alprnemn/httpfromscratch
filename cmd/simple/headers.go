package main

import "bytes"

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

	if err := validateHeaders(headers); err != nil {
		return 0, false, err
	}

	for i := 0; i < len(headers); i++ {
		parts := bytes.SplitN(headers[i], []byte(":"), 2)
		h[string(parts[0])] = string(parts[1])
	}

	return read, true, nil
}

// TODO: Validate headers ' ' by RFC make whitespaces appropriate
func validateHeaders(b [][]byte) error {
	return nil
}

//POST /coffee HTTP/1.1\r\n
//    Host: localhost:42069     \r\n
//User-Agent: curl/7.81.0\r\n
//Accept: */*\r\n
//Content-Length: 11\r\n
//\r\n
//hello world

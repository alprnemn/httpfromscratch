package response

import (
	"errors"
	"fmt"
	"httpfromscratch/common"
	headers "httpfromscratch/internal/headers"
	"net"
)

type StatusCode int
type StatusMsg string

const (
	StatusOK                  StatusCode = 200
	StatusBadRequestError     StatusCode = 400
	StatusInternalServerError StatusCode = 500
	StatusOKMsg               StatusMsg  = "HTTP/1.1 200 OK\r\n"
	StatusBRMsg               StatusMsg  = "HTTP/1.1 400 Bad Request\r\n"
	StatusISEMsg              StatusMsg  = "HTTP/1.1 500 Internal Server Error\r\n"
)

type Writer struct {
	conn    net.Conn
	Headers *headers.Headers
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		conn:    conn,
		Headers: GetDefaultHeaders(),
	}
}

// WriteChunkedBody writes a single HTTP/1.1 chunk to the underlying connection.
// It formats the data according to the chunked transfer encoding specification:
//
//	<chunk-size in hex>\r\n
//	<chunk-data>\r\n
//
// If the provided byte slice is empty, the function performs no operation.
func (w *Writer) WriteChunkedBody(p []byte) error {

	if len(p) == 0 {
		return nil
	}

	hexStr, _ := common.ConvertDecToHex(len(p))

	if _, err := w.Write([]byte(hexStr + "\r\n")); err != nil {
		return err
	}

	if _, err := w.Write(p); err != nil {
		return err
	}

	if _, err := w.Write([]byte("\r\n")); err != nil {
		return err
	}

	return nil
}

// WriteChunkedBodyDone writes the terminating chunk that signals the end
// of a chunked HTTP/1.1 response body. According to the specification,
// this consists of a zero-length chunk:
//
//	0\r\n\r\n
func (w *Writer) WriteChunkedBodyDone() error {
	_, err := w.Write([]byte("0\r\n\r\n"))
	return err
}

// WriteStatusLine writes the HTTP status line to the connection.
// The status line includes the protocol version, status code,
// and reason phrase, and must be written before any headers or body data.
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.conn.Write([]byte(StatusOKMsg))
		if err != nil {
			return err
		}
	case StatusBadRequestError:
		_, err := w.conn.Write([]byte(StatusBRMsg))
		if err != nil {
			return err
		}

	case StatusInternalServerError:
		_, err := w.conn.Write([]byte(StatusISEMsg))
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteHeaders writes all HTTP response headers to the connection,
// followed by the mandatory blank line that separates headers from the body.
func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	headers.ForEach(func(n, v string) {
		line := fmt.Sprintf("%s: %s\r\n", n, v)
		_, err := w.conn.Write([]byte(line))
		if err != nil {
			return
		}
	})

	_, err := w.conn.Write([]byte("\r\n"))
	return err
}

// Write writes raw bytes directly to the underlying TCP connection.
// If the payload is non-empty, it updates the Content-Length header
// to reflect the size of the written data.
//
// This method should not be used when Transfer-Encoding is set to "chunked".
func (w *Writer) Write(p []byte) (int, error) {

	_, err := w.conn.Write(p)
	if err != nil {
		return 0, errors.New("error occurred while writin body")
	}
	return 0, nil
}

// GetDefaultHeaders returns a default set of HTTP response headers
// configured for a chunked response body.
func GetDefaultHeaders() *headers.Headers {
	h := headers.NewHeaders()
	h.SetHeader("Content-Type", "text/plain")
	h.SetHeader("Transfer-Encoding", "chunked")
	return h
}

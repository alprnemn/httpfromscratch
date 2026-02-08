package response

import (
	"errors"
	"fmt"
	h "httpfromscratch/internal/headers"
	"net"
	"strconv"
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
	headers *h.Headers
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		conn:    conn,
		headers: GetDefaultHeaders(),
	}
}

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

func (w *Writer) WriteHeaders(headers *h.Headers) error {
	headers.ForEach(func(n, v string) {
		line := fmt.Sprintf("%s: %s\r\n", n, v)
		w.conn.Write([]byte(line))
	})

	_, err := w.conn.Write([]byte("\r\n"))
	return err
}

func (w *Writer) Write(p []byte) (int, error) {

	if len(p) > 0 {
		strLength := strconv.Itoa(len(p))
		w.headers.Set("Content-length", strLength)
	}

	_, err := w.conn.Write(p)
	if err != nil {
		return 0, errors.New("error occurred while writin body")
	}
	return 0, nil
}

func GetDefaultHeaders() *h.Headers {
	headers := h.NewHeaders()
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

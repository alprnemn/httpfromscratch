package response

import (
	"errors"
	"fmt"
	h "httpfromscratch/internal/headers"
	"io"
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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	switch statusCode {
	case StatusOK:
		if err := writeStatus(w, StatusOKMsg); err != nil {
			return err
		}
	case StatusBadRequestError:
		if err := writeStatus(w, StatusBRMsg); err != nil {
			return err
		}
	case StatusInternalServerError:
		if err := writeStatus(w, StatusISEMsg); err != nil {
			return err
		}
	}
	return nil
}

func writeStatus(w io.Writer, msg StatusMsg) error {

	headers := GetDefaultHeaders(5)

	// write response msg
	_, err := w.Write([]byte(msg))
	if err != nil {
		return errors.New("error writing status ok")
	}

	// write headers
	if err := WriteHeaders(w, headers); err != nil {
		return err
	}

	// write body
	_, err = w.Write([]byte("\r\nHello"))
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) *h.Headers {
	headers := h.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers *h.Headers) error {
	headers.ForEach(func(n, v string) {
		hdr := fmt.Sprintf("%s: %s\r\n", n, v)
		_, err := w.Write([]byte(hdr))
		if err != nil {
			return
		}
	})
	return nil
}

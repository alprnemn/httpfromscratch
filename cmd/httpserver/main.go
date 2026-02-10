package main

import (
	headers "httpfromscratch/internal/headers"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"httpfromscratch/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const port = 42069
const binURL = "https://httpbin.org/stream/10"
const BufferSizeChunk = 1024

func main() {

	sv, err := server.Serve(port, proxyHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer sv.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		if err := WriteResponse(w, response.StatusBadRequestError, "this is your fault!!"); err != nil {
			log.Fatal(err)
		}
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		if err := WriteResponse(w, response.StatusInternalServerError, "server's fault!!"); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := WriteResponse(w, response.StatusOK, "Server is running succesfull!!"); err != nil {
		log.Fatal(err)
	}
	return

}

func proxyHandler(w *response.Writer, req *request.Request) {
	res, err := http.Get(binURL)
	if err != nil {
		err := WriteResponse(w, response.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
	defer res.Body.Close()

	if err := w.WriteStatusLine(response.StatusOK); err != nil {
		err := WriteResponse(w, response.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}

	w.WriteHeaders(w.Headers)

	if err := streamChunked(w, res.Body); err != nil {
		log.Printf("stream error: %v", err)
	}

}

// streamChunked reads data from the provided io.Reader and streams it
// to the client using HTTP/1.1 chunked transfer encoding.
//
// It continuously reads into a fixed-size buffer and writes each
// chunk using Writer.WriteChunkedBody. When the reader signals EOF,
// the terminating chunk is written using Writer.WriteChunkedBodyDone.
func streamChunked(w *response.Writer, reader io.Reader) error {
	buf := make([]byte, BufferSizeChunk)
	for {
		time.Sleep(time.Second)
		n, err := reader.Read(buf)
		if n > 0 {
			if err := w.WriteChunkedBody(buf[:n]); err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				return w.WriteChunkedBodyDone()
			}
			return err
		}

	}
}

// WriteResponse writes a complete HTTP response using the provided
// status code and message body. It writes the status line, headers,
// and body in the correct HTTP/1.1 order.
func WriteResponse(w *response.Writer, statusCode response.StatusCode, msg string) error {
	if err := w.WriteStatusLine(statusCode); err != nil {
		return err
	}

	h := headers.NewHeaders()
	h.SetHeader("Content-Length", strconv.Itoa(len(msg)))

	err := w.WriteHeaders(h)
	if err != nil {
		return err
	}

	if err := WriteBodyMsg(w, msg); err != nil {
		return err
	}

	return nil
}

// WriteBodyMsg writes the given message string directly to the response body
// using the provided response.Writer.
func WriteBodyMsg(w *response.Writer, msg string) error {
	_, err := w.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

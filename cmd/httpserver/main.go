package main

import (
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"httpfromscratch/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func WriteBodyMsg(w *response.Writer, msg string) error {
	_, err := w.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func WriteResponse(w *response.Writer, statusCode response.StatusCode, msg string) error {
	if err := w.WriteStatusLine(statusCode); err != nil {
		return err
	}

	headers := response.GetDefaultHeaders()

	err := w.WriteHeaders(headers)
	if err != nil {
		return err
	}
	if err := WriteBodyMsg(w, msg); err != nil {
		return err
	}
	return nil
}

func main() {

	handler := func(w *response.Writer, req *request.Request) {

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

		if err := WriteResponse(w, response.StatusOK, "success"); err != nil {
			log.Fatal(err)
		}
		return

	}

	sv, err := server.Serve(port, handler)

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

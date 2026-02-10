package server

import (
	"httpfromscratch/internal/request"
	response "httpfromscratch/internal/response"
)

// Handler defines the function signature for processing an incoming
// HTTP request and generating a response using the provided Writer.
type Handler func(w *response.Writer, req *request.Request)

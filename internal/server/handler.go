package server

import (
	"httpfromscratch/internal/request"
	r "httpfromscratch/internal/response"
)

type HandlerError struct {
	Code    r.StatusCode
	Message string
}

type Handler func(w *r.Writer, req *request.Request)

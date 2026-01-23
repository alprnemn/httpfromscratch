package main

import "errors"

type ReqState string

const BufferSize = 8

const (
	StateDone    ReqState = "done"
	StateInit    ReqState = "init"
	StateHeaders ReqState = "headers"
	StateError   ReqState = "error"
)

type Request struct {
	RequestLine RequestLine
	Headers     Headers
	Body        []byte
	State       ReqState
}

func NewRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (r Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

// request parse
func (r *Request) Parse(data []byte) (int, error) {

	read := 0

outer:
	for {
		switch r.State {
		case StateError:
			return 0, errors.New("error state error")
		case StateInit:
			rl, readedFromParse, err := ParseRequestLine(data)
			if err != nil {
				r.State = StateError
				return 0, err
			}

			if readedFromParse == 0 {
				break outer
			}

			read += readedFromParse

			r.RequestLine = *rl
			r.State = StateHeaders
		}

	}

	return read, nil
}

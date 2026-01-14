package main

import (
	"errors"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	State       parserState
}

func NewRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0

outer:
	for {
		switch r.State {
		case StateError:
			return 0, errors.New("error in request state")
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])

			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.State = StateDone

		case StateDone:
			break outer
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

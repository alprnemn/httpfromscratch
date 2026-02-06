# HTTP Server From Scratch (Go)

A minimal **HTTP/1.1 server implemented from scratch in Go** using raw TCP sockets.  
This project focuses on understanding HTTP internals by rebuilding core server components without using `net/http`.


## Features

- Raw TCP-based HTTP server
- HTTP/1.1 request parsing
- Custom response writer abstraction
- Flexible handler system
- Header & status code management
- Concurrent connection handling
- HTML response support

## Purpose

To deeply understand how HTTP servers work internally by implementing:

- TCP networking
- HTTP request parsing
- Response serialization
- Header management
- Concurrency with goroutines

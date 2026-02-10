.PHONY: push udpsender udpreceiver tcplistener

# run listeners

run sv:
	@go run ./cmd/httpserver

udpl:
	@go run ./cmd/udpsender

tcpl:
	@go run ./cmd/tcplistener

udpreceiver:
	nc -u -l 42069

# tests
test-request:
	@go test ./internal/request
test-headers:
	@go test ./internal/headers
test-helpers:
	@go test ./common

ccd:
	@go run ./cmd/ccd
simple:
	@go run ./internal/simple

get:
	@curl -v --request GET -sL \
	     --url 'http://localhost:42069/httpbin/stream/100' \
	     -H "Accept: application/json" \
  		 -H "Authorization: Bearer <your_token_here>"



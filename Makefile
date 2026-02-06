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

# tests but raise warning
test request:
	@go test ./internal/request
test headers:
	@go test ./internal/headers

ccd:
	@go run ./cmd/ccd
simple:
	@go run ./internal/simple

get:
	@curl -v --request GET -sL \
	     --url 'http://localhost:42069/health' \
	     -H "Accept: application/json" \
  		 -H "Authorization: Bearer <your_token_here>"



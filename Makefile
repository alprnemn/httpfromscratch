.PHONY: push udpsender udpreceiver tcplistener

# run listeners

udpl:
	@go run ./cmd/udpsender
tcpl:
	@go run ./cmd/tcplistener

udpreceiver:
	nc -u -l 42069

# tests
test request:
	@go test ./internal/request
test headers:
	@go test ./internal/headers


simple:
	@go run ./cmd/simple

get:
	@curl --request GET -sL \
	     --url 'http://localhost:42069/coffee' \
	     -H "Accept: application/json" \
  		 -H "Authorization: Bearer your_token_here" \
  		 -H "Authori zation: Bearer asdfasfdasf"



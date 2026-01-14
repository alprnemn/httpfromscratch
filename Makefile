.PHONY: push udpsender udpreceiver tcplistener

udpsender:
	@go run ./cmd/udpsender

tcplistener:
	@go run ./cmd/tcplistener

udpreceiver:
	nc -u -l 42069

test request:
	@go test ./internal/request

simple:
	@go run ./cmd/simple

get:
	@curl --request GET -sL \
	     --url 'http://localhost:42069/coffee'

push:
ifndef MSG
	$(error MSG is required. Usage: make push MSG="commit message")
endif
	@git add .
	sleep 1
	@git commit -m "$(MSG)"
	sleep 2
	@git push origin main

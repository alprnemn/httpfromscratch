udpsender:
	@go run ./cmd/udpsender

tcplistener:
	@go run ./cmd/tcplistener

udpreceiver:
	nc -u -l 42069
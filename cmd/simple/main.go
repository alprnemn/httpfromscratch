package main

import (
	"fmt"
	h "httpfromscratch/internal/headers"
	"log"
)

func main() {

	headers := h.NewHeaders()

	//data := []byte("Host:   localhost:42069        \r\n      User-Agent:      curl/7.81.0       \r\n    Accept: */*   \r\nContent-Length: 11\r\n\r\nhello world")
	data := []byte("       Host: localhost:42069\r\n      Host:      localhost:55555        \nUser-Agent: Curl\n\n\n\n")

	_, _, err := headers.ParseHeader(data)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range headers {
		fmt.Printf("--%s--,--%s--\n", k, v)
	}

}

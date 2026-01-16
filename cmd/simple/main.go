package main

import (
	"fmt"
	h "httpfromscratch/internal/headers"
	"log"
)

func main() {

	headers := h.NewHeaders()

	data := []byte("Host:   localhost:42069        \r\n      User-Agent:      curl/7.81.0       \r\n    Accept: */*   \r\nContent-Length: 11\r\n\r\nhello world")

	_, _, err := headers.ParseHeader(data)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range headers {
		fmt.Printf("key:%s\n", k)
		fmt.Printf("value:%s\n", v)
	}

}

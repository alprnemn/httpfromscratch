package main

type Headers struct {
	Headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		Headers: map[string]string{},
	}
}







// GET /products HTTP/1.1\r\nHost: localhost:44433\r\nUser-Agent: Curl.7.2\r\n\r\nFooBaar

package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaders_Parse(t *testing.T) {
	h := NewHeaders()

	data := []byte("Host: localhost:42069\r\nHost: localhost:43069\r\n\r\n")

	_, done, err := h.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069,localhost:43069", h.Get("Host"))
	//assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	//assert.Equal(t, "*/*", headers["accept"])
	//assert.Equal(t, "11", headers["content-length"])
	assert.True(t, done)

	//headers = NewHeaders()
	//data = []byte("H©st: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Length: 11\r\n\r\nhello world")
	//n, done, err := headers.Parse(data)
	//require.Error(t, err)
	//assert.Equal(t, 0, n)
	//assert.Equal(t, "invalid header key: H©st", err.Error())
	//assert.False(t, done)
	//
	//headers = NewHeaders()
	//data = []byte("       Host: loca\r\nlhost:42069       \r\n\r\n")
	//n, done, err = headers.Parse(data)
	//assert.Equal(t, "loca", headers["host"])
	//
	//headers = NewHeaders()
	//data = []byte("       Host: localhost:42069\n      Host:      localhost:55555        \nUser-Agent: Curl\r\n\r\n")
	//n, done, err = headers.Parse(data)
	//assert.Equal(t, "localhost:42069,localhost:55555", headers["host"])

}

//POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\n
//User-Agent: curl/7.81.0\r\n
//Accept: */*\r\n
//Content-Length: 11\r\n
//\r\n
//hello world

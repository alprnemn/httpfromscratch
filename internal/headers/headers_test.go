package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaders_Parse(t *testing.T) {
	headers := NewHeaders()

	data := []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Length: 11\r\n\r\nhello world")

	_, done, err := headers.ParseHeader(data)

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, "*/*", headers["accept"])
	assert.Equal(t, "11", headers["content-length"])
	assert.True(t, done)

	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Length: 11\r\n\r\nhello world")
	n, done, err := headers.ParseHeader(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, "invalid header key: H©st", err.Error())
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("       Host: loca\r\nlhost:42069       \r\n\r\n")
	n, done, err = headers.ParseHeader(data)
	assert.Equal(t, "loca", headers["host"])

	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\n      Host:      localhost:55555        \nUser-Agent: Curl\r\n\r\n")
	n, done, err = headers.ParseHeader(data)
	assert.Equal(t, "localhost:42069,localhost:55555", headers["host"])

}

//POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\n
//User-Agent: curl/7.81.0\r\n
//Accept: */*\r\n
//Content-Length: 11\r\n
//\r\n
//hello world

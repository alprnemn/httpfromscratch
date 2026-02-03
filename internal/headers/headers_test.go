package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaders_Parse(t *testing.T) {
	h := NewHeaders()
	data := []byte("User-Agent: curl/7.81.0\r\nHost: localhost:42069\r\nHost: localhost:43069\r\n\r\n")
	_, done, err := h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	host, ok := h.Get("Host")
	assert.Equal(t, "localhost:42069,localhost:43069", host)
	assert.Equal(t, true, ok)
	useragent, ok := h.Get("User-Agent")
	assert.Equal(t, "curl/7.81.0", useragent)
	assert.Equal(t, true, ok)
	assert.True(t, done)

	data = []byte("H©st: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nContent-Length: 11\r\n\r\nhello world")
	n, done, err := h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, "invalid header key: H©st", err.Error())
	assert.False(t, done)

	data = []byte("       Host: loca\r\nlhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "42069", h.headers["lhost"])

	headers := NewHeaders()
	data = []byte("       Host: localhost:42069\r\n      Host:      localhost:     55555        \r\nUser-Agent: Curl\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	host, ok = headers.Get("host")
	assert.Equal(t, "localhost:42069,localhost:     55555", host)
	assert.True(t, ok)

	head := NewHeaders()
	data = []byte("       Host: localhost:42069\n      Host:      localhost:     55555        \nUser-Agent: Curl\r\n\r\n")
	n, done, err = head.Parse(data)
	require.Error(t, err)
	assert.Equal(t, "invalid header value: localhost:42069\n      Host:      localhost:     55555        \nUser-Agent: Curl", err.Error())

	header1 := NewHeaders()
	data = []byte("Host: localhost:42069\r\n      Host1:      localhost:     55555        \r\nUser-Agent: Curl\r\nContent-Length: 11\r\n\r\n")
	n, done, err = header1.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", header1.headers["host"])
}

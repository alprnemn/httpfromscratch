package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParseRequestLine(t *testing.T) {
	r, err := ParseRequestFromReader(strings.NewReader(
		"POST /products HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/products", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	_, err = ParseRequestFromReader(strings.NewReader(
		"PoST /products HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	))
	require.Error(t, err)
	assert.Equal(t, "method chars must be capitalize", err.Error())

	r, err = ParseRequestFromReader(strings.NewReader(
		"POST /products HTTP/3.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	))
	require.Error(t, err)
	assert.Equal(t, "only supported version is 1.1", err.Error())

	r, err = ParseRequestFromReader(strings.NewReader(
		"POST %products HTTP/3.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	))
	require.Error(t, err)
	assert.Equal(t, "target must start /", err.Error())

	r, err = ParseRequestFromReader(strings.NewReader(
		"POST /products HTTP/3.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	))
	require.Error(t, err)
	assert.Equal(t, "only supported version is 1.1", err.Error())

	r, err = ParseRequestFromReader(strings.NewReader(
		"/products HTTP/3.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	))
	require.Error(t, err)
	assert.Equal(t, "request line parts must be 3", err.Error())

}

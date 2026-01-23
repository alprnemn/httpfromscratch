package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseRequestLine(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /asd HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := ParseRequest(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/asd", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "localhost:42069", r.Headers.Get("host"))
	assert.Equal(t, "curl/7.81.0", r.Headers.Get("user-agent"))

	reader = &chunkReader{
		data:            "GET /asd HTTP/1.1\r\nHos t: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*",
		numBytesPerRead: 3,
	}
	r, err = ParseRequest(reader)
	require.Error(t, err)
	assert.Equal(t, "invalid header key: Hos t", err.Error())

	reader = &chunkReader{
		data:            "GET /asd HTTP/1.1\r\nH©st:: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = ParseRequest(reader)
	require.Error(t, err)
	assert.Equal(t, "invalid header key: H©st", err.Error())

	reader = &chunkReader{
		data:            "GET /asd HTTP/1.1\r\nHost:: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = ParseRequest(reader)
	//require.Error(t, err)
	assert.Equal(t, ": localhost:42069", r.Headers.Get("host"))

	reader = &chunkReader{
		data:            "ASD /asd HTTP/1.1\r\nHost:: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 5,
	}
	r, err = ParseRequest(reader)
	require.Error(t, err)
	assert.Equal(t, "method name not implemented", err.Error())

}

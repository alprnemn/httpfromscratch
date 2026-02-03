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

	host, ok := r.Headers.Get("host")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, true, ok)

	useragent, ok := r.Headers.Get("user-agent")
	assert.Equal(t, true, ok)
	assert.Equal(t, "curl/7.81.0", useragent)

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
	host, ok = r.Headers.Get("host")
	assert.Equal(t, ": localhost:42069", host)

	reader = &chunkReader{
		data:            "ASD /asd HTTP/1.1\r\nHost:: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 5,
	}
	r, err = ParseRequest(reader)
	require.Error(t, err)
	assert.Equal(t, "method name not implemented", err.Error())

	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
		numBytesPerRead: 3,
	}
	r, err = ParseRequest(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", string(r.Body))

	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	}
	r, err = ParseRequest(reader)
	require.NoError(t, err)
	host, ok = r.Headers.Get("host")
	assert.Equal(t, "localhost:42069", host)

	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length:0\r\n" +
			"\r\n" +
			"asd",
		numBytesPerRead: 3,
	}
	r, err = ParseRequest(reader)
	require.NoError(t, err)
	host, ok = r.Headers.Get("host")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, "asd", string(r.Body))

	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length:11\r\n" +
			"\r\n" +
			"Hello World",
		numBytesPerRead: 3,
	}
	r, err = ParseRequest(reader)
	require.NoError(t, err)
	host, ok = r.Headers.Get("host")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, "asd", string(r.Body))

}

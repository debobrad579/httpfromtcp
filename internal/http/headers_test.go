package http_test

import (
	"testing"

	"github.com/debobrad579/httpfromtcp/internal/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	h := http.NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h.Get("host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	h = http.NewHeaders()
	data = []byte("       Host: localhost:42069                           \r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h.Get("host"))
	assert.Equal(t, 57, n)
	assert.False(t, done)

	// Test: Valid 2 h with existing h
	h = http.NewHeaders()
	h.Set("host", "localhost:42069")
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h.Get("host"))
	assert.Equal(t, "curl/7.81.0", h.Get("user-agent"))
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid 2 h with same field name
	h = http.NewHeaders()
	h.Set("host", "localhost:42069")
	data = []byte("Host: localhost:8080\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069, localhost:8080", h.Get("host"))
	assert.Equal(t, 22, n)
	assert.False(t, done)

	// Test: Valid done
	h = http.NewHeaders()
	data = []byte("\r\n{'flavor':'dark mode'}")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid numbers in field name
	h = http.NewHeaders()
	data = []byte("H0st123456789: localhost:42069\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h.Get("h0st123456789"))
	assert.Equal(t, 32, n)
	assert.False(t, done)

	// Test: Valid special characters in field name
	h = http.NewHeaders()
	data = []byte("Ho$t!#%&'*+-.^_`|~: localhost:42069\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h.Get("ho$t!#%&'*+-.^_`|~"))
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Invalid special characters in field name
	h = http.NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid blank field name
	h = http.NewHeaders()
	data = []byte(": localhost:42069\r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing between field name and colon
	h = http.NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

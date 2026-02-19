package http

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	i := bytes.Index(data, []byte("\r\n"))
	if i == -1 {
		return nil, 0, nil
	}

	requestLine := string(data[:i])
	consumed := i + 2
	parts := strings.Fields(requestLine)

	if len(parts) != 3 {
		return nil, 0, errors.New("invalid number of parts in request line")
	}

	method := parts[0]
	for _, letter := range method {
		if !(unicode.IsLetter(letter) && unicode.IsUpper(letter)) {
			return nil, 0, fmt.Errorf("incorrect method format: %s", method)
		}
	}

	reqTarget := parts[1]

	if reqTarget == "" || !strings.HasPrefix(reqTarget, "/") {
		return nil, 0, fmt.Errorf("invalid request target format")
	}

	if !strings.HasPrefix(parts[2], "HTTP/") {
		return nil, 0, fmt.Errorf("invalid http version format: %s", parts[2])
	}
	httpVersion := strings.TrimPrefix(parts[2], "HTTP/")
	if httpVersion != "1.1" {
		return nil, 0, fmt.Errorf("incorrect http version: %s", httpVersion)
	}

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: reqTarget,
		Method:        method,
	}, consumed, nil
}

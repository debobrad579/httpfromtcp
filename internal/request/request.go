package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const bufferSize = 8

type requestState int

const (
	requestInitialized requestState = iota
	requestDone
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	req := &Request{state: requestInitialized}

	for req.state != requestDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		nRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				req.state = requestDone
				break
			}
			return nil, err
		}

		readToIndex += nRead

		nParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[nParsed:readToIndex])
		readToIndex -= nParsed
	}

	return req, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	i := strings.Index(string(data), "\r\n")
	if i == -1 {
		return RequestLine{}, 0, nil
	}

	reqLine := string(data[:i])
	consumed := i + 2
	parts := strings.Split(reqLine, " ")

	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("Invalid number of parts in request line")
	}

	method := parts[0]
	for _, letter := range method {
		if !unicode.IsLetter(letter) || !unicode.IsUpper(letter) {
			return RequestLine{}, 0, fmt.Errorf("Incorrect method format: %s", method)
		}
	}

	reqTarget := parts[1]

	httpVersion := strings.TrimPrefix(parts[2], "HTTP/")
	if httpVersion != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("Incorrect http version: %s", httpVersion)
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: reqTarget,
		Method:        method,
	}, consumed, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestInitialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = reqLine
		r.state = requestDone
		return n, nil
	case requestDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

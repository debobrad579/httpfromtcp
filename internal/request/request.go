package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/debobrad579/httpfromtcp/internal/headers"
)

const bufferSize = 8

type requestState int

const (
	requestInitialized requestState = iota
	requestDone
	requestParsingHeaders
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       requestState
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	request := &Request{state: requestInitialized}
	request.Headers = make(headers.Headers)

	for request.state != requestDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		nRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			// if err == io.EOF {
			//	 request.state = requestDone
			//	 break
			// }
			return nil, err
		}

		readToIndex += nRead

		nParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[nParsed:readToIndex])
		readToIndex -= nParsed
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != requestDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return n, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}

		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = requestLine
		r.state = requestParsingHeaders
		return n, nil
	case requestParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.state = requestDone
		}

		return n, nil
	case requestDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	i := strings.Index(string(data), "\r\n")
	if i == -1 {
		return RequestLine{}, 0, nil
	}

	requestLine := string(data[:i])
	consumed := i + 2
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("error: invalid number of parts in request line")
	}

	method := parts[0]
	for _, letter := range method {
		if !unicode.IsLetter(letter) || !unicode.IsUpper(letter) {
			return RequestLine{}, 0, fmt.Errorf("error: incorrect method format: %s", method)
		}
	}

	reqTarget := parts[1]

	httpVersion := strings.TrimPrefix(parts[2], "HTTP/")
	if httpVersion != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("error: incorrect http version: %s", httpVersion)
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: reqTarget,
		Method:        method,
	}, consumed, nil
}

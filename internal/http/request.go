package http

import (
	"errors"
	"io"
	"strconv"
)

const bufferSize = 1024
const maxBufferSize = 8 * 1024 * 1024

type requestState int

const (
	requestInitialized requestState = iota
	requestParsingHeaders
	requestParsingBody
	requestDone
)

type Request struct {
	RequestLine RequestLine
	Headers     Headers
	Body        []byte
	state       requestState
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	request := &Request{state: requestInitialized}
	request.Headers = *NewHeaders()

	for request.state != requestDone {
		if readToIndex >= len(buf) {
			if len(buf)*2 > maxBufferSize {
				return nil, errors.New("request too large")
			}
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		nRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
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

		r.RequestLine = *requestLine
		r.state = requestParsingHeaders
		return n, nil
	case requestParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.state = requestParsingBody
		}

		return n, nil
	case requestParsingBody:
		contentLengthStr := r.Headers.Get("Content-Length")
		if contentLengthStr == "" {
			r.state = requestDone
			return 0, nil
		}

		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, err
		}

		if contentLength == 0 {
			r.state = requestDone
			return 0, nil
		}

		remaining := contentLength - len(r.Body)
		if len(data) > remaining {
			data = data[:remaining]
		}

		r.Body = append(r.Body, data...)

		if len(r.Body) == contentLength {
			r.state = requestDone
		}

		return len(data), nil
	case requestDone:
		return 0, errors.New("trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
}

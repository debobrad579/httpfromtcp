package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return parseRequestLine(string(reqData))
}

func parseRequestLine(reqString string) (*Request, error) {
	req_line := strings.Split(reqString, "\r\n")[0]
	parts := strings.Split(req_line, " ")

	if len(parts) != 3 {
		return nil, errors.New("Invalid number of parts in request line")
	}

	method := parts[0]
	for _, letter := range method {
		if !unicode.IsLetter(letter) || !unicode.IsUpper(letter) {
			return nil, fmt.Errorf("Incorrect method format: %s", method)
		}
	}

	reqTarget := parts[1]

	httpVersion := strings.TrimPrefix(parts[2], "HTTP/")
	if httpVersion != "1.1" {
		return nil, fmt.Errorf("Incorrect http version: %s", httpVersion)
	}

	req := Request{
		RequestLine: RequestLine{
			HttpVersion:   httpVersion,
			RequestTarget: reqTarget,
			Method:        method,
		},
	}

	return &req, nil
}

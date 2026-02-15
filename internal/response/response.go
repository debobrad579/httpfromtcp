package response

import (
	"fmt"
	"net"
	"strconv"

	"github.com/debobrad579/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Writer struct {
	Conn net.Conn
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.Conn.Write(p)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var reasonPhrase string

	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		reasonPhrase = "Unknown"
	}

	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for fieldName, fieldValue := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", fieldName, fieldValue)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w, "\r\n")
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Write(p)
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := make(headers.Headers)
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

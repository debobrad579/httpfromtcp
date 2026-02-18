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
	StatusNotFound            StatusCode = 404
	StatusMethodNotAllowed    StatusCode = 405
	StatusInternalServerError StatusCode = 500
)

var reasonPhrases = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusInternalServerError: "Internal Server Error",
}

type Writer struct {
	conn net.Conn
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{conn}
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.conn.Write(p)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	reasonPhrase, ok := reasonPhrases[statusCode]
	if !ok {
		reasonPhrase = "Unknown"
	}

	if _, err := fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase); err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := make(headers.Headers)
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	for fieldName, fieldValue := range h {
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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	return fmt.Fprintf(w, "%x\r\n%s\r\n", len(p), p)
}

func (w *Writer) WriteChunkedBodyDone(trailers headers.Headers) error {
	if _, err := fmt.Fprint(w, "0\r\n"); err != nil {
		return err
	}

	for fieldName, fieldValue := range trailers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", fieldName, fieldValue); err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(w, "\r\n")
	return err
}

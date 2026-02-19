package http

import (
	"fmt"
	"io"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type ResponseWriter struct {
	writerState writerState
	writer      io.Writer
}

func NewResponseWriter(writer io.Writer) *ResponseWriter {
	return &ResponseWriter{
		writerState: writerStateStatusLine,
		writer:      writer,
	}
}

func (w *ResponseWriter) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateHeaders }()

	reasonPhrase, ok := reasonPhrases[statusCode]
	if !ok {
		reasonPhrase = "Unknown"
	}

	if _, err := fmt.Fprintf(w.writer, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase); err != nil {
		return err
	}

	return nil
}

func (w *ResponseWriter) WriteHeaders(h *Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateBody }()

	var writeErr error
	h.Range(func(fieldName, fieldValue string) bool {
		if _, writeErr = fmt.Fprintf(w.writer, "%s: %s\r\n", fieldName, fieldValue); writeErr != nil {
			return false
		}
		return true
	})

	if writeErr != nil {
		return writeErr
	}

	_, err := fmt.Fprint(w.writer, "\r\n")
	return err
}

func (w *ResponseWriter) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	return w.writer.Write(p)
}

func (w *ResponseWriter) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write chunked body in state %d", w.writerState)
	}

	return fmt.Fprintf(w.writer, "%x\r\n%s\r\n", len(p), p)
}

func (w *ResponseWriter) WriteTrailers(trailers *Headers) error {
	if w.writerState != writerStateBody {
		return fmt.Errorf("cannot write trailers in state %d", w.writerState)
	}

	if _, err := fmt.Fprint(w.writer, "0\r\n"); err != nil {
		return err
	}

	var writeErr error

	trailers.Range(func(fieldName, fieldValue string) bool {
		if _, writeErr = fmt.Fprintf(w.writer, "%s: %s\r\n", fieldName, fieldValue); writeErr != nil {
			return false
		}
		return true
	})

	if writeErr != nil {
		return writeErr
	}

	_, err := fmt.Fprint(w.writer, "\r\n")
	return err
}

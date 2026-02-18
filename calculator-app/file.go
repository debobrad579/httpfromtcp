package main

import (
	"mime"
	"os"
	"path/filepath"

	"github.com/debobrad579/httpfromtcp/internal/response"
)

func writeFileResponse(w *response.Writer, path string) {
	html, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteStatusLine(response.StatusNotFound)
		} else {
			w.WriteStatusLine(response.StatusInternalServerError)
		}
		return
	}

	w.WriteStatusLine(200)
	headers := response.GetDefaultHeaders(len(html))
	headers.Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
	w.WriteHeaders(headers)
	w.WriteBody(html)
}

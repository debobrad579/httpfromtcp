package main

import (
	"mime"
	"os"
	"path/filepath"

	"github.com/debobrad579/httpfromtcp/internal/http"
)

func writeFileResponse(w *http.ResponseWriter, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteStatusLine(http.StatusNotFound)
		} else {
			w.WriteStatusLine(http.StatusInternalServerError)
		}
		return
	}

	w.WriteStatusLine(http.StatusOK)
	headers := http.GetDefaultResponseHeaders(mime.TypeByExtension(filepath.Ext(path)), len(data))
	w.WriteHeaders(headers)
	w.WriteBody(data)
}

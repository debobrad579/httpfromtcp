package main

import (
	"strings"

	"github.com/debobrad579/httpfromtcp/internal/request"
	"github.com/debobrad579/httpfromtcp/internal/response"
)

func routeHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/" {
		writeFileResponse(w, "calculator-app/templates/index.html")
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/static/") {
		writeFileResponse(w, "calculator-app"+req.RequestLine.RequestTarget)
		return
	}

	if req.RequestLine.RequestTarget == "/api" {
		apiHandler(w, req)
		return
	}

	w.WriteStatusLine(response.StatusNotFound)
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/debobrad579/httpfromtcp/internal/headers"
	"github.com/debobrad579/httpfromtcp/internal/request"
	"github.com/debobrad579/httpfromtcp/internal/response"
	"github.com/debobrad579/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, routeHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func routeHandler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		httpBinProxyHandler(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {
	default:
		notFoundHandler(w, req)
	}
}

func httpBinProxyHandler(w *response.Writer, req *request.Request) {
	x := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	res, err := http.Get("https://httpbin.org/" + x)
	if err != nil {
		internalServerErrorHandler(w, req)
		return
	}
	defer res.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := make(headers.Headers)
	h.Set("Content-Type", "application/json")
	h.Set("Connection", "close")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailers", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(h)

	var fullBody []byte
	buf := make([]byte, 1024)

	for {
		n, err := res.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			return
		}

		fullBody = append(fullBody, buf[:n]...)
		w.WriteChunkedBody(buf[:n])
	}

	hash := sha256.Sum256(fullBody)
	hashString := hex.EncodeToString(hash[:])

	trailers := make(headers.Headers)
	trailers.Set("X-Content-SHA256", hashString)
	trailers.Set("X-Content-Length", strconv.Itoa(len(fullBody)))

	w.WriteChunkedBodyDone(trailers)
}

func writeHTMLResponse(w *response.Writer, html string, statusCode response.StatusCode) {
	w.WriteStatusLine(statusCode)
	headers := response.GetDefaultHeaders(len(html))
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(html))
}

func notFoundHandler(w *response.Writer, _ *request.Request) {
	html := `<html>
		<head>
			<title>404 - Not Found</title>
		</head>
		<body>
			<h1>404 - Not Found</h1>
		</body>
	</html>`

	writeHTMLResponse(w, html, response.StatusNotFound)
}

func internalServerErrorHandler(w *response.Writer, _ *request.Request) {
	html := `<html>
		<head>
			<title>500 - Internal Server Error</title>
		</head>
		<body>
			<h1>500 - Internal Server Error</h1>
		</body>
	</html>`

	writeHTMLResponse(w, html, response.StatusInternalServerError)
}

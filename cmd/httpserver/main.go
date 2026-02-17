package main

import (
	"encoding/json"
	"log"
	"mime"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/Pramod-Devireddy/go-exprtk"
	"github.com/debobrad579/httpfromtcp/internal/request"
	"github.com/debobrad579/httpfromtcp/internal/response"
	"github.com/debobrad579/httpfromtcp/internal/server"
)

const port = 8080

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
	if req.RequestLine.RequestTarget == "/" {
		writeFileResponse(w, "templates/index.html")
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/static/") {
		writeFileResponse(w, req.RequestLine.RequestTarget[1:])
		return
	}

	if req.RequestLine.RequestTarget == "/api" {
		apiHandler(w, req)
		return
	}

	notFound(w)
}

func writeFileResponse(w *response.Writer, path string) {
	html, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			notFound(w)
		} else {
			internalServerError(w, err.Error())
		}
		return
	}

	w.WriteStatusLine(200)
	headers := response.GetDefaultHeaders(len(html))
	headers.Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
	w.WriteHeaders(headers)
	w.WriteBody(html)
}

func notFound(w *response.Writer) {
	w.WriteStatusLine(404)
	w.WriteHeaders(response.GetDefaultHeaders(0))
}

func internalServerError(w *response.Writer, message string) {
	w.WriteStatusLine(500)
	w.WriteHeaders(response.GetDefaultHeaders(len(message)))
	w.WriteBody([]byte(message))
}

type apiResponseBody struct {
	Equation string `json:"equation"`
	IsError  bool   `json:"is_error"`
}

type apiRequestBody struct {
	Equation     string `json:"equation"`
	IsDegreeMode bool   `json:"is_degree_mode"`
}

func apiHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.Method != "POST" {
		return
	}

	var reqBody apiRequestBody
	if err := json.Unmarshal(req.Body, &reqBody); err != nil {
		return
	}

	if reqBody.IsDegreeMode {
		reqBody.Equation = regexp.MustCompile(`(^|[^a])sin\(([^)]+)\)`).ReplaceAllString(reqBody.Equation, "${1}sin(pi/180*($2))")
		reqBody.Equation = regexp.MustCompile(`(^|[^a])cos\(([^)]+)\)`).ReplaceAllString(reqBody.Equation, "${1}cos(pi/180*($2))")
		reqBody.Equation = regexp.MustCompile(`(^|[^a])tan\(([^)]+)\)`).ReplaceAllString(reqBody.Equation, "${1}tan(pi/180*($2))")
		reqBody.Equation = regexp.MustCompile(`asin\(([^)]+)\)`).ReplaceAllString(reqBody.Equation, "(180/pi*asin($1))")
		reqBody.Equation = regexp.MustCompile(`acos\(([^)]+)\)`).ReplaceAllString(reqBody.Equation, "(180/pi*acos($1))")
		reqBody.Equation = regexp.MustCompile(`atan\(([^)]+)\)`).ReplaceAllString(reqBody.Equation, "(180/pi*atan($1))")
	}

	expr := exprtk.NewExprtk()
	expr.SetExpression(reqBody.Equation)
	expr.CompileExpression()

	resBody := apiResponseBody{
		Equation: strings.TrimSuffix(strings.TrimRight(strconv.FormatFloat(expr.GetEvaluatedValue(), 'f', 9, 64), "0"), "."),
		IsError:  false,
	}

	resData, err := json.Marshal(resBody)
	if err != nil {
		return
	}

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(len(resData))
	h.Set("Content-Type", "application/json")
	w.WriteHeaders(h)
	w.WriteBody(resData)
}

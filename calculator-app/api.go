package main

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/Pramod-Devireddy/go-exprtk"
	"github.com/debobrad579/httpfromtcp/internal/request"
	"github.com/debobrad579/httpfromtcp/internal/response"
)

type apiResponseBody struct {
	EvaluatedValue string `json:"evaluated_value"`
}

type apiRequestBody struct {
	Equation     string `json:"equation"`
	IsDegreeMode bool   `json:"is_degree_mode"`
}

func apiHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.Method != "POST" {
		w.WriteStatusLine(response.StatusMethodNotAllowed)
		return
	}

	var reqBody apiRequestBody
	if err := json.Unmarshal(req.Body, &reqBody); err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
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
	if err := expr.CompileExpression(); err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		return
	}

	resBody := apiResponseBody{
		EvaluatedValue: strconv.FormatFloat(expr.GetEvaluatedValue(), 'g', 9, 64),
	}

	resData, err := json.Marshal(resBody)
	if err != nil {
		w.WriteStatusLine(response.StatusInternalServerError)
		return
	}

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(len(resData))
	h.Set("Content-Type", "application/json")
	w.WriteHeaders(h)
	w.WriteBody(resData)
}

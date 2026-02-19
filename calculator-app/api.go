package main

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/Pramod-Devireddy/go-exprtk"
	"github.com/debobrad579/httpfromtcp/internal/http"
)

type apiResponseBody struct {
	EvaluatedValue string `json:"evaluated_value"`
}

type apiRequestBody struct {
	Equation     string `json:"equation"`
	IsDegreeMode bool   `json:"is_degree_mode"`
}

func apiHandler(w *http.ResponseWriter, req *http.Request) {
	if req.RequestLine.Method != "POST" {
		w.WriteStatusLine(http.StatusMethodNotAllowed)
		return
	}

	var reqBody apiRequestBody
	if err := json.Unmarshal(req.Body, &reqBody); err != nil {
		w.WriteStatusLine(http.StatusBadRequest)
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
		w.WriteStatusLine(http.StatusBadRequest)
		return
	}

	evaluatedValue := expr.GetEvaluatedValue()
	if evaluatedValue < 0.00000001 && evaluatedValue > -0.00000001 {
		evaluatedValue = 0.0
	}

	resBody := apiResponseBody{
		EvaluatedValue: strconv.FormatFloat(evaluatedValue, 'g', 9, 64),
	}

	resData, err := json.Marshal(resBody)
	if err != nil {
		w.WriteStatusLine(http.StatusInternalServerError)
		return
	}

	w.WriteStatusLine(http.StatusOK)
	h := http.GetDefaultResponseHeaders("application/json", len(resData))
	w.WriteHeaders(h)
	w.WriteBody(resData)
}

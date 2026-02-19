package http

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	h := Headers{}
	h.headers = make(map[string]string)
	return &h
}

func GetDefaultResponseHeaders(mimetype string, contentLen int) *Headers {
	h := NewHeaders()
	h.Set("Connection", "close")
	h.Set("Content-Type", mimetype)
	h.Set("Content-Length", strconv.Itoa(contentLen))
	return h
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h *Headers) Set(key, value string) {
	h.headers[strings.ToLower(key)] = value
}

func (h *Headers) Del(key string) {
	delete(h.headers, strings.ToLower(key))
}

func (h *Headers) Range(callback func(key, value string) bool) {
	for k, v := range h.headers {
		if !callback(k, v) {
			return
		}
	}
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	i := bytes.Index(data, []byte("\r\n"))
	if i == -1 {
		return 0, false, nil
	}

	fieldLine := data[:i]

	if len(fieldLine) == 0 {
		return 2, true, nil
	}

	fieldName, fieldValue, found := strings.Cut(string(fieldLine), ":")
	if !found || strings.HasSuffix(fieldName, " ") {
		return 0, false, errors.New("invalid header format: " + fieldName)
	}

	fieldName = strings.ToLower(strings.TrimSpace(fieldName))
	if !isValidFieldName(fieldName) {
		return 0, false, errors.New("invalid header format: " + fieldName)
	}

	fieldValue = strings.TrimSpace(fieldValue)
	prevFieldValue := h.Get(fieldName)
	if prevFieldValue != "" {
		fieldValue = prevFieldValue + ", " + fieldValue
	}

	h.Set(fieldName, fieldValue)
	return i + 2, false, nil
}

func isValidFieldName(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+.^_`|~-", c) {
			return false
		}
	}

	return true
}

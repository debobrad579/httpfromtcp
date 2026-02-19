package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers struct {
	h map[string]string
}

func New() *Headers {
	headers := Headers{}
	headers.h = make(map[string]string)
	return &headers
}

func (h *Headers) Get(key string) string {
	return h.h[strings.ToLower(key)]
}

func (h *Headers) Set(key, value string) {
	h.h[strings.ToLower(key)] = value
}

func (h *Headers) Del(key string) {
	delete(h.h, strings.ToLower(key))
}

func (h *Headers) Range(callback func(key, value string) bool) {
	for k, v := range h.h {
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

package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func (h Headers) Get(fieldName string) string {
	return h[strings.ToLower(fieldName)]
}

func (h Headers) Set(fieldName string, fieldValue string) {
	h[strings.ToLower(fieldName)] = fieldValue
}

func (h Headers) Del(fieldName string) {
	delete(h, strings.ToLower(fieldName))
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	i := bytes.Index(data, []byte("\r\n"))
	if i == -1 {
		return 0, false, nil
	}

	line := data[:i]

	if len(line) == 0 {
		return 2, true, nil
	}

	fieldName, fieldValue, found := strings.Cut(string(line), ":")
	if !found || strings.HasSuffix(fieldName, " ") {
		return 0, false, errors.New("error: invalid header format: " + fieldName)
	}

	fieldName = strings.ToLower(strings.TrimSpace(fieldName))
	fieldNamePattern := "^[a-z0-9!#$%&'*+.^_`|~-]+$"
	matched, err := regexp.MatchString(fieldNamePattern, fieldName)
	if fieldName == "" || !matched {
		return 0, false, errors.New("error: invalid header format: " + fieldName)
	}

	fieldValue = strings.TrimSpace(fieldValue)
	prevFieldValue, ok := h[fieldName]
	if ok {
		fieldValue = prevFieldValue + ", " + fieldValue
	}

	h[fieldName] = fieldValue
	return i + 2, false, nil
}

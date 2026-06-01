package utils

import (
	"io"
	"net/http"
)

func NewTestRequest(method, target string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, target, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}
